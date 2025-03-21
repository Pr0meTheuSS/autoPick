package repository

import (
	"context"
	"errors"
	"user-service/internal/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository - репозиторий для работы с пользователями
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository - конструктор для репозитория пользователей
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// Create - создание нового пользователя
func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	user.ID = uuid.NewString()
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID - получение пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("неверный формат ID")
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("пользователь не найден")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("пользователь не найден")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

// List - получение списка пользователей с пагинацией
func (r *UserRepository) List(ctx context.Context, page, limit int) ([]models.User, error) {
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit64)

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Update - обновление данных пользователя
func (r *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	update := bson.M{
		"$set": bson.M{
			"username": user.Username,
			"password": user.Password,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, errors.New("пользователь не найден")
	}

	return r.GetByID(ctx, user.ID)
}

// Delete - удаление пользователя по ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("неверный формат ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("пользователь не найден")
	}

	return nil
}

func (r *UserRepository) GetRefreshTokenByUserId(ctx context.Context, userID string) (string, error) {
	var refresh models.RefreshToken
	if err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&refresh); err != nil {
		return "", err
	}
	return refresh.RefreshToken, nil
}

func (r *UserRepository) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) (*models.RefreshToken, error) {
	_, err := r.collection.InsertOne(ctx, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

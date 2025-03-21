package repository

import (
	"context"
	"product-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductRepository - репозиторий для работы с продуктами в MongoDB
type ProductRepository struct {
	collection *mongo.Collection
}

// NewProductRepository - конструктор для создания нового репозитория
func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

// Create - создание нового продукта в базе данных
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

// GetByID - получение продукта по ID
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// List - получение списка всех продуктов
// List - получение списка продуктов
func (r *ProductRepository) List(ctx context.Context) ([]models.Product, error) {
	var products []models.Product

	// Используем options.Find() для возможности кастомизации запроса в будущем
	findOptions := options.Find()

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// Проверяем ошибку после окончания итерации по курсору
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// Update - обновление существующего продукта в базе данных
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	// Мы используем ID продукта как уникальный идентификатор для поиска
	filter := bson.M{"_id": product.ID}

	// Обновляем только измененные поля
	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
		},
	}

	// Выполняем операцию обновления
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete - удаление продукта по ID
func (r *ProductRepository) Delete(ctx context.Context, product *models.Product) error {
	// Удаляем продукт по ID
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": product.ID})
	return err
}

package repository

import (
	"context"
	"order-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OrderRepository - репозиторий для работы с заказами в MongoDB
type OrderRepository struct {
	collection *mongo.Collection
}

// NewOrderRepository - конструктор для создания нового репозитория
func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		collection: db.Collection("orders"),
	}
}

// Create - создание нового заказа в базе данных
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// GetByID - получение заказа по ID
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// List - получение списка заказов
func (r *OrderRepository) List(ctx context.Context, filter bson.M, limit, offset int) ([]models.Order, error) {
	var orders []models.Order

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.M{"created_at": -1}) // Сортировка по дате создания (новые сначала)

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// Update - обновление существующего заказа
func (r *OrderRepository) Update(ctx context.Context, order *models.Order) (*models.Order, error) {
	filter := bson.M{"_id": order.ID}
	update := bson.M{"$set": order}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return order, err
}

// Delete - удаление заказа по ID
func (r *OrderRepository) Delete(ctx context.Context, orderID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": orderID})
	return err
}

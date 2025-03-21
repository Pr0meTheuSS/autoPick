package usecase

import (
	"context"
	"errors"
	"order-service/internal/models"
	"order-service/internal/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// OrderService - сервис для работы с продуктами
type OrderService struct {
	repo *repository.OrderRepository
}

// NewOrderService - конструктор для создания сервиса
func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// Create - создание нового продукта
func (s *OrderService) Create(ctx context.Context, Order *models.Order) (*models.Order, error) {
	// Генерация уникального ID для продукта
	Order.ID = uuid.NewString()
	return s.repo.Create(ctx, Order)
}

// GetByID - получение продукта по ID
func (s *OrderService) GetByID(ctx context.Context, id string) (*models.Order, error) {
	return s.repo.GetByID(ctx, id)
}

// List - получение списка всех продуктов
func (s *OrderService) List(ctx context.Context, filter bson.M, limit, offset int) ([]models.Order, error) {
	return s.repo.List(ctx, filter, limit, offset)
}

// Update - обновление данных продукта
func (s *OrderService) UpdateOrderStatus(ctx context.Context, Order *models.Order) (*models.Order, error) {
	// Проверка, существует ли продукт с таким ID
	existingOrder, err := s.repo.GetByID(ctx, Order.ID)
	if err != nil {
		return nil, errors.New("продукт не найден")
	}

	// Обновление данных продукта
	existingOrder.Status = Order.Status
	// Сохранение обновленного продукта в базе
	return s.repo.Update(ctx, existingOrder)
}

// Delete - удаление продукта по ID
func (s *OrderService) Delete(ctx context.Context, id string) error {
	// Проверка, существует ли продукт с таким ID
	Order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("продукт не найден")
	}

	// Удаление продукта из базы
	return s.repo.Delete(ctx, Order.ID)
}

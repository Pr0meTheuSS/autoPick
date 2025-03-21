package usecase

import (
	"context"
	"errors"
	"product-service/internal/models"
	"product-service/internal/repository"

	"github.com/google/uuid"
)

// ProductService - сервис для работы с продуктами
type ProductService struct {
	repo *repository.ProductRepository
}

// NewProductService - конструктор для создания сервиса
func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// Create - создание нового продукта
func (s *ProductService) Create(ctx context.Context, product *models.Product) error {
	// Генерация уникального ID для продукта
	product.ID = uuid.NewString()
	return s.repo.Create(ctx, product)
}

// GetByID - получение продукта по ID
func (s *ProductService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// List - получение списка всех продуктов
func (s *ProductService) List(ctx context.Context) ([]models.Product, error) {
	return s.repo.List(ctx)
}

// Update - обновление данных продукта
func (s *ProductService) Update(ctx context.Context, product *models.Product) error {
	// Проверка, существует ли продукт с таким ID
	existingProduct, err := s.repo.GetByID(ctx, product.ID)
	if err != nil {
		return errors.New("продукт не найден")
	}

	// Обновление данных продукта
	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price

	// Сохранение обновленного продукта в базе
	return s.repo.Update(ctx, existingProduct)
}

// Delete - удаление продукта по ID
func (s *ProductService) Delete(ctx context.Context, id string) error {
	// Проверка, существует ли продукт с таким ID
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("продукт не найден")
	}

	// Удаление продукта из базы
	return s.repo.Delete(ctx, product)
}

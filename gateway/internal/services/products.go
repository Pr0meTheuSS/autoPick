package services

import (
	"context"
	"fmt"

	"gateway/internal/models"
	"gateway/internal/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProductsService - gRPC клиент
type ProductsService struct {
	client proto.ProductServiceClient
	logger *zap.Logger
}

// NewProductsService - конструктор сервиса с логированием
func NewProductsService(grpcAddress string, logger *zap.Logger) (*ProductsService, error) {
	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Ошибка подключения к gRPC серверу", zap.String("address", grpcAddress), zap.Error(err))
		return nil, fmt.Errorf("не удалось подключиться к gRPC серверу: %w", err)
	}

	client := proto.NewProductServiceClient(conn)
	logger.Info("gRPC клиент успешно подключен", zap.String("address", grpcAddress))

	return &ProductsService{client: client, logger: logger}, nil
}

// Get - получение списка продуктов с логированием
func (p *ProductsService) Get(ctx context.Context, page int, limit int) ([]models.Product, error) {
	p.logger.Info("Запрос списка продуктов", zap.Int("page", page), zap.Int("limit", limit))

	resp, err := p.client.GetProducts(ctx, &proto.GetProductsRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		p.logger.Error("Ошибка получения списка продуктов", zap.Error(err))
		return nil, err
	}

	var products []models.Product
	for _, p := range resp.Products {
		products = append(products, models.Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	p.logger.Info("Список продуктов успешно получен", zap.Int("count", len(products)))
	return products, nil
}

// GetByID - получение продукта по ID с логированием
func (p *ProductsService) GetByID(ctx context.Context, id string) (models.Product, error) {
	p.logger.Info("Запрос продукта по ID", zap.String("id", id))

	resp, err := p.client.GetProductByID(ctx, &proto.GetProductByIDRequest{Id: id})
	if err != nil {
		return models.Product{}, err
	}

	product := models.Product{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		Price:       resp.Price,
	}

	p.logger.Info("Продукт успешно получен", zap.String("id", product.ID))
	return product, nil
}

// Create - создание нового продукта с логированием
func (p *ProductsService) Create(product models.Product) (models.Product, error) {
	p.logger.Info("Создание нового продукта", zap.String("name", product.Name))

	resp, err := p.client.CreateProduct(context.Background(), &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       float32(product.Price),
	})
	if err != nil {
		p.logger.Error("Ошибка создания продукта", zap.String("name", product.Name), zap.Error(err))
		return models.Product{}, err
	}

	createdProduct := models.Product{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		Price:       resp.Price,
	}

	p.logger.Info("Продукт успешно создан", zap.String("id", createdProduct.ID))
	return createdProduct, nil
}

// Put - обновление продукта с логированием
func (p *ProductsService) Put(product models.Product) (models.Product, error) {
	p.logger.Info("Обновление продукта", zap.String("id", product.ID))

	resp, err := p.client.UpdateProduct(context.Background(), &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	})
	if err != nil {
		p.logger.Error("Ошибка обновления продукта", zap.String("id", product.ID), zap.Error(err))
		return models.Product{}, err
	}

	updatedProduct := models.Product{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		Price:       resp.Price,
	}

	p.logger.Info("Продукт успешно обновлен", zap.String("id", updatedProduct.ID))
	return updatedProduct, nil
}

// Delete - удаление продукта по ID с логированием
func (p *ProductsService) Delete(id string) error {
	p.logger.Info("Удаление продукта", zap.String("id", id))

	_, err := p.client.DeleteProduct(context.Background(), &proto.DeleteProductRequest{Id: id})
	if err != nil {
		p.logger.Error("Ошибка удаления продукта", zap.String("id", id), zap.Error(err))
		return err
	}

	p.logger.Info("Продукт успешно удален", zap.String("id", id))
	return nil
}

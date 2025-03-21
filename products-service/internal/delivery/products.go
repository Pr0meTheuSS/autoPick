package delivery

import (
	"context"
	"product-service/internal/models"
	"product-service/internal/proto"
	"product-service/internal/usecase"

	"go.uber.org/zap" // Импортируем zap для логирования
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ proto.ProductServiceServer = (*ProductHandler)(nil)

// ProductHandler - обработчик запросов для продуктов
type ProductHandler struct {
	proto.UnimplementedProductServiceServer
	service *usecase.ProductService
	logger  *zap.Logger // Логгер
}

// NewProductHandler - конструктор для создания обработчика
func NewProductHandler(service *usecase.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{service: service, logger: logger}
}

// CreateProduct - обработка запроса на создание продукта
func (h *ProductHandler) CreateProduct(ctx context.Context, req *proto.Product) (*proto.Product, error) {
	h.logger.Info("Получен запрос на создание продукта", zap.String("name", req.Name))

	// Создание модели продукта из запроса
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	// Вызов бизнес-логики для создания продукта
	err := h.service.Create(ctx, product)
	if err != nil {
		// Логирование ошибки и возврат ошибки с соответствующим кодом
		h.logger.Error("Ошибка при создании продукта", zap.String("name", req.Name), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось создать продукт: %v", err)
	}

	// Логирование успешного создания
	h.logger.Info("Продукт успешно создан", zap.String("id", product.ID), zap.String("name", product.Name))

	// Возвращаем ответ с созданным продуктом
	return &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

// GetProductByID - обработка запроса на получение продукта по ID
func (h *ProductHandler) GetProductByID(ctx context.Context, req *proto.GetProductByIDRequest) (*proto.Product, error) {
	h.logger.Info("Получен запрос на получение продукта", zap.String("id", req.Id))

	// Получаем продукт по ID
	product, err := h.service.GetByID(ctx, req.Id)
	if err != nil {
		// Логирование ошибки и возврат ошибки с кодом NotFound, если продукт не найден
		h.logger.Error("Продукт не найден", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "не удалось получить продукт с ID %s: %v", req.Id, err)
	}

	// Логирование успешного получения продукта
	h.logger.Info("Продукт успешно найден", zap.String("id", product.ID))

	// Возвращаем ответ с найденным продуктом
	return &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

// GetProducts - обработка запроса на получение списка продуктов
func (h *ProductHandler) GetProducts(ctx context.Context, req *proto.GetProductsRequest) (*proto.GetProductsResponse, error) {
	h.logger.Info("Получен запрос на получение списка продуктов")

	// Получаем все продукты
	products, err := h.service.List(ctx)
	if err != nil {
		// Логирование ошибки и возврат ошибки с соответствующим кодом
		h.logger.Error("Ошибка при получении списка продуктов", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось получить список продуктов: %v", err)
	}

	// Логирование успешного получения списка продуктов
	h.logger.Info("Список продуктов успешно получен", zap.Int("count", len(products)))

	// Формируем ответ с продуктами
	var productList []*proto.Product
	for _, product := range products {
		productList = append(productList, &proto.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	// Возвращаем ответ с продуктами
	return &proto.GetProductsResponse{
		Products: productList,
	}, nil
}

// UpdateProduct - обработка запроса на обновление продукта
func (h *ProductHandler) UpdateProduct(ctx context.Context, req *proto.Product) (*proto.Product, error) {
	h.logger.Info("Получен запрос на обновление продукта", zap.String("id", req.Id))

	// Создание модели продукта из запроса
	product := &models.Product{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	// Вызов бизнес-логики для обновления продукта
	err := h.service.Update(ctx, product)
	if err != nil {
		// Логирование ошибки и возврат ошибки с соответствующим кодом
		h.logger.Error("Ошибка при обновлении продукта", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось обновить продукт: %v", err)
	}

	// Логирование успешного обновления
	h.logger.Info("Продукт успешно обновлен", zap.String("id", req.Id))

	// Возвращаем ответ с обновленным продуктом
	return &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

// DeleteProduct - обработка запроса на удаление продукта
func (h *ProductHandler) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	h.logger.Info("Получен запрос на удаление продукта", zap.String("id", req.Id))

	// Вызов бизнес-логики для удаления продукта
	err := h.service.Delete(ctx, req.Id)
	if err != nil {
		// Логирование ошибки и возврат ошибки с соответствующим кодом
		h.logger.Error("Ошибка при удалении продукта", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось удалить продукт с ID %s: %v", req.Id, err)
	}

	// Логирование успешного удаления
	h.logger.Info("Продукт успешно удален", zap.String("id", req.Id))

	// Возвращаем ответ с подтверждением удаления
	return &proto.DeleteProductResponse{
		Success: true,
	}, nil
}

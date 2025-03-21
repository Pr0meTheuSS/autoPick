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

// OrdersService - gRPC клиент
type OrdersService struct {
	client proto.OrderServiceClient
	logger *zap.Logger
}

// NewOrdersService - конструктор сервиса с логированием
func NewOrdersService(grpcAddress string, logger *zap.Logger) (*OrdersService, error) {
	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Ошибка подключения к gRPC серверу", zap.String("address", grpcAddress), zap.Error(err))
		return nil, fmt.Errorf("не удалось подключиться к gRPC серверу: %w", err)
	}

	client := proto.NewOrderServiceClient(conn)
	logger.Info("gRPC клиент успешно подключен", zap.String("address", grpcAddress))

	return &OrdersService{client: client, logger: logger}, nil
}

// Get - получение списка заказов
func (o *OrdersService) Get(ctx context.Context, offset int, limit int, userID string) ([]models.Order, error) {
	o.logger.Info("Запрос списка заказов", zap.Int("page", offset), zap.Int("limit", limit))

	resp, err := o.client.ListOrders(ctx, &proto.ListOrdersRequest{
		UserId: userID,
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		o.logger.Error("Ошибка получения списка заказов", zap.Error(err))
		return nil, err
	}

	var orders []models.Order
	for _, ord := range resp.Orders {
		orders = append(orders, models.Order{
			ID:          ord.Id,
			UserID:      ord.UserId,
			ProductsIDs: ord.ProductIds,
			Status:      ord.Status,
			Total:       uint(ord.TotalPrice),
			CreatedAt:   ord.GetCreatedAt().AsTime(),
		})
	}

	o.logger.Info("Список заказов успешно получен", zap.Int("count", len(orders)))
	return orders, nil
}

// GetByID - получение заказа по ID
func (o *OrdersService) GetByID(ctx context.Context, id string) (models.Order, error) {
	o.logger.Info("Запрос заказа по ID", zap.String("id", id))

	resp, err := o.client.GetOrder(ctx, &proto.GetOrderRequest{
		OrderId: id,
	})
	if err != nil {
		return models.Order{}, err
	}

	order := models.Order{
		ID:          resp.Order.Id,
		UserID:      resp.Order.UserId,
		ProductsIDs: resp.Order.ProductIds,
		Status:      resp.Order.Status,
		Total:       uint(resp.Order.TotalPrice),
	}

	o.logger.Info("Заказ успешно получен", zap.String("id", order.ID))
	return order, nil
}

// Create - создание нового заказа
func (o *OrdersService) Create(order models.Order) (models.Order, error) {
	o.logger.Info("Создание нового заказа", zap.String("user_id", order.UserID))

	resp, err := o.client.CreateOrder(context.Background(), &proto.CreateOrderRequest{
		UserId:     order.UserID,
		ProductIds: order.ProductsIDs,
		Total:      int32(order.Total),
	})
	if err != nil {
		o.logger.Error("Ошибка создания заказа", zap.String("user_id", order.UserID), zap.Error(err))
		return models.Order{}, err
	}

	createdOrder := models.Order{
		ID:          resp.GetOrder().GetId(),
		UserID:      resp.GetOrder().GetUserId(),
		ProductsIDs: resp.GetOrder().GetProductIds(),
		Status:      resp.GetOrder().GetStatus(),
		Total:       uint(resp.GetOrder().GetTotalPrice()),
		CreatedAt:   resp.GetOrder().GetCreatedAt().AsTime(),
	}

	o.logger.Info("Заказ успешно создан", zap.String("id", createdOrder.ID))
	return createdOrder, nil
}

// Update - обновление заказа
func (o *OrdersService) Update(order models.Order) (models.Order, error) {
	o.logger.Info("Обновление заказа", zap.String("id", order.ID))

	_, err := o.client.UpdateOrderStatus(context.Background(), &proto.UpdateOrderStatusRequest{
		OrderId: order.ID,
		Status:  order.Status,
	})
	if err != nil {
		o.logger.Error("Ошибка обновления заказа", zap.String("id", order.ID), zap.Error(err))
		return models.Order{}, err
	}

	o.logger.Info("Заказ успешно обновлен", zap.String("id", order.ID))
	return order, nil
}

// Delete - удаление заказа
func (o *OrdersService) Delete(id string) error {
	o.logger.Info("Удаление заказа", zap.String("id", id))

	_, err := o.client.DeleteOrder(context.Background(), &proto.DeleteOrderRequest{
		OrderId: id,
	})

	if err != nil {
		o.logger.Error("Ошибка удаления заказа", zap.String("id", id), zap.Error(err))
		return err
	}

	o.logger.Info("Заказ успешно удален", zap.String("id", id))
	return nil
}

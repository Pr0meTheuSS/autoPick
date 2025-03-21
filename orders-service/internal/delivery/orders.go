package delivery

import (
	"context"
	"order-service/internal/models"
	"order-service/internal/proto"
	"order-service/internal/usecase"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ proto.OrderServiceServer = (*OrderHandler)(nil)

type OrderHandler struct {
	proto.UnimplementedOrderServiceServer
	service *usecase.OrderService
	logger  *zap.Logger
}

func NewOrderHandler(service *usecase.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{service: service, logger: logger}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	h.logger.Info("Создание заказа", zap.String("user_id", req.UserId))

	order := &models.Order{
		UserID:     req.UserId,
		ProductIDs: req.ProductIds,
		Status:     "pending",
	}

	savedOrder, err := h.service.Create(ctx, order)
	if err != nil {
		h.logger.Error("Ошибка создания заказа", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось создать заказ: %v", err)
	}

	return &proto.CreateOrderResponse{
		Order: &proto.Order{
			Id:         savedOrder.ID,
			UserId:     savedOrder.UserID,
			ProductIds: savedOrder.ProductIDs,
			TotalPrice: savedOrder.TotalPrice,
			Status:     savedOrder.Status,
			CreatedAt:  timestamppb.New(savedOrder.CreatedAt),
		},
	}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	h.logger.Info("Получение списка заказов", zap.String("user_id", req.UserId))

	filter := bson.M{
		"user_id": req.UserId,
	}
	limit := req.Limit
	offset := req.Offset

	orders, err := h.service.List(ctx, filter, int(limit), int(offset))
	if err != nil {
		h.logger.Error("Ошибка при получении заказов", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось получить заказы: %v", err)
	}

	var protoOrders []*proto.Order
	for _, order := range orders {
		protoOrders = append(protoOrders, convertToProtoOrder(&order))
	}

	return &proto.ListOrdersResponse{Orders: protoOrders}, nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *proto.UpdateOrderStatusRequest) (*proto.UpdateOrderStatusResponse, error) {
	h.logger.Info("Обновление статуса заказа", zap.String("order_id", req.OrderId), zap.String("status", req.Status))

	order := &models.Order{
		ID:     req.OrderId,
		Status: req.Status,
	}

	_, err := h.service.UpdateOrderStatus(ctx, order)
	if err != nil {
		h.logger.Error("Ошибка обновления статуса", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось обновить статус: %v", err)
	}

	return &proto.UpdateOrderStatusResponse{Success: true}, nil
}

func (h *OrderHandler) DeleteOrder(ctx context.Context, req *proto.DeleteOrderRequest) (*proto.DeleteOrderResponse, error) {
	h.logger.Info("Удаление заказа", zap.String("order_id", req.OrderId))

	if err := h.service.Delete(ctx, req.OrderId); err != nil {
		h.logger.Error("Ошибка удаления заказа", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось удалить заказ: %v", err)
	}

	return &proto.DeleteOrderResponse{Success: true}, nil
}

func convertToProtoOrder(order *models.Order) *proto.Order {
	return &proto.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		ProductIds: order.ProductIDs,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  timestamppb.New(order.CreatedAt),
	}
}

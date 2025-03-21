package handlers

import (
	"encoding/json"
	"fmt"
	"gateway/internal/dtos"
	"gateway/internal/middleware"
	"gateway/internal/models"
	"gateway/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// OrdersHandler - обработчик заказов
type OrdersHandler struct {
	service services.OrdersService
}

// NewOrdersHandler - конструктор обработчика заказов
func NewOrdersHandler(service services.OrdersService) *OrdersHandler {
	return &OrdersHandler{
		service: service,
	}
}

// GetOrders godoc
// @Summary Получить список заказов
// @Description Возвращает список заказов с возможностью пагинации
// @Tags orders
// @Accept  json
// @Produce  json
// @Param offset query int false "offset" default(0)
// @Param limit query int false "Items per page" default(10)
// @Param user_id query string false "Уникальный идентификатор пользователя, который офрмлял заказы"
// @Success 200 {array} dtos.OrderDto
// @Failure 500 {string} string "Ошибка сервера"
// @Router /orders [get]
func (o *OrdersHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paginationParams, ok := middleware.GetPaginationParamsFromCtx(r.Context())
	if !ok {
		paginationParams = &middleware.PaginationParams{
			Limit:  10,
			Offset: 0,
		}
	}

	userID := r.URL.Query().Get("user_id")

	fmt.Println(userID)
	orders, err := o.service.Get(ctx, paginationParams.Offset, paginationParams.Limit, userID)
	if err != nil {
		http.Error(w, "Не удалось получить заказы", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetOrderByID godoc
// @Summary Получить заказ по ID
// @Description Возвращает информацию о заказе по его ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "ID заказа"
// @Success 200 {object} dtos.OrderDto
// @Failure 400 {string} string "Некорректный ID"
// @Failure 404 {string} string "Заказ не найден"
// @Router /orders/{id} [get]
func (o *OrdersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Не передан id заказа", http.StatusBadRequest)
		return
	}

	order, err := o.service.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// CreateOrder godoc
// @Summary Создать заказ
// @Description Добавляет новый заказ
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body dtos.CreateOrderDto true "Данные нового заказа"
// @Success 201 {object} dtos.OrderDto
// @Failure 400 {string} string "Неверные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /orders [post]
func (o *OrdersHandler) Post(w http.ResponseWriter, r *http.Request) {
	var dto dtos.CreateOrderDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	order := models.Order{
		UserID:      dto.UserID,
		ProductsIDs: dto.ProductsIDs,
		Total:       dto.Total,
	}

	createdOrder, err := o.service.Create(order)
	if err != nil {
		http.Error(w, "Ошибка при создании заказа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdOrder)
}

// UpdateOrder godoc
// @Summary Обновить заказ
// @Description Обновляет существующий заказ
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body dtos.OrderDto true "Обновленные данные заказа"
// @Success 200 {object} dtos.OrderDto
// @Failure 400 {string} string "Неверные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /orders [put]
func (o *OrdersHandler) Put(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	updatedOrder, err := o.service.Update(order)
	if err != nil {
		http.Error(w, "Ошибка при обновлении заказа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

// DeleteOrder godoc
// @Summary Удалить заказ
// @Description Удаляет заказ по ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "ID заказа"
// @Success 204 "Заказ удален"
// @Failure 400 {string} string "Некорректный ID"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /orders/{id} [delete]
func (o *OrdersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Не передан id заказа", http.StatusBadRequest)
		return
	}

	if err := o.service.Delete(id); err != nil {
		http.Error(w, "Ошибка при удалении заказа", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

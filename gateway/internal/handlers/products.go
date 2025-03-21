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

// ProductsHandler - обработчик продуктов
type ProductsHandler struct {
	service services.ProductsService
}

// NewProductsHandler - конструктор обработчика продуктов
func NewProductsHandler(service services.ProductsService) *ProductsHandler {
	return &ProductsHandler{
		service: service,
	}
}

// GetProducts godoc
// @Summary Получить список продуктов
// @Description Возвращает список продуктов с возможностью пагинации
// @Tags products
// @Accept  json
// @Produce  json
// @Param offset query int false "offset" default(0)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} dtos.ProductDto
// @Failure 500 {string} string "Ошибка сервера"
// @Router /products [get]
func (p *ProductsHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paginationParams, ok := middleware.GetPaginationParamsFromCtx(r.Context())
	if !ok {
		paginationParams = &middleware.PaginationParams{
			Limit:  10,
			Offset: 0,
		}
	}

	// Вызываем сервис
	products, err := p.service.Get(ctx, paginationParams.Offset, paginationParams.Limit)
	if err != nil {
		http.Error(w, "Не удалось получить продукты", http.StatusInternalServerError)
		return
	}

	// Преобразуем в JSON и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetProductByID godoc
// @Summary Получить продукт по ID
// @Description Возвращает информацию о продукте по его ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "ID продукта"
// @Success 200 {object} dtos.ProductDto
// @Failure 400 {string} string "Некорректный ID"
// @Failure 404 {string} string "Продукт не найден"
// @Router /products/{id} [get]
func (p *ProductsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	fmt.Println(id)
	fmt.Println(r.URL)
	if id == "" {
		http.Error(w, "Не передан id продукта", http.StatusBadRequest)
		return
	}

	product, err := p.service.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Продукт не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// CreateProduct godoc
// @Summary Создать продукт
// @Description Добавляет новый продукт
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body dtos.CreateProductDto true "Данные нового продукта"
// @Success 201 {object} dtos.ProductDto
// @Failure 400 {string} string "Неверные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /products [post]
func (p *ProductsHandler) Post(w http.ResponseWriter, r *http.Request) {
	var dto dtos.CreateProductDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Преобразуем DTO в модель
	product := models.Product{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
	}

	createdProduct, err := p.service.Create(product)
	if err != nil {
		http.Error(w, "Ошибка при создании продукта", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

// UpdateProduct godoc
// @Summary Обновить продукт
// @Description Обновляет существующий продукт
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body dtos.ProductDto true "Обновленные данные продукта"
// @Success 200 {object} dtos.ProductDto
// @Failure 400 {string} string "Неверные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /products [put]
func (p *ProductsHandler) Put(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	updatedProduct, err := p.service.Put(product)
	if err != nil {
		http.Error(w, "Ошибка при обновлении продукта", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

// DeleteProduct godoc
// @Summary Удалить продукт
// @Description Удаляет продукт по ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "ID продукта"
// @Success 204 "Продукт удален"
// @Failure 400 {string} string "Некорректный ID"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /products/{id} [delete]
func (p *ProductsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Не передан id продукта", http.StatusBadRequest)
		return
	}

	if err := p.service.Delete(id); err != nil {
		http.Error(w, "Ошибка при удалении продукта", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

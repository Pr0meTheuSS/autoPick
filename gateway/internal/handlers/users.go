package handlers

import (
	"encoding/json"
	"gateway/internal/dtos"
	"gateway/internal/middleware"
	"gateway/internal/models"
	"gateway/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// UserHandler - обработчик запросов для пользователей
type UserHandler struct {
	service *services.UsersService
}

// NewUserHandler - конструктор для создания обработчика
func NewUserHandler(service *services.UsersService) *UserHandler {
	return &UserHandler{service: service}
}

// GetUsers godoc
// @Summary Получить список пользователей с пагинацией
// @Description Получение списка пользователей с возможностью пагинации по страницам и лимиту.
// @Tags users
// @Accept  json
// @Produce  json
// @Param offset query int false "offset" default(0)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} dtos.UserDto
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Server error"
// @Router /users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры пагинации из запроса
	paginationParams, ok := middleware.GetPaginationParamsFromCtx(r.Context())
	if !ok {
		paginationParams = &middleware.PaginationParams{
			Limit:  10,
			Offset: 0,
		}
	}

	// Запрашиваем список пользователей с пагинацией
	users, err := h.service.GetUsers(r.Context(), paginationParams.Offset, paginationParams.Limit)
	if err != nil {
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}

	// Отправляем результат
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByID godoc
// @Summary Получить пользователя по ID
// @Description Получение данных пользователя по ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} dtos.UserDto
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser godoc
// @Summary Создать нового пользователя
// @Description Создание нового пользователя с предоставленными данными. Статус почты - неподтвержденный, роль - пользователь.
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body dtos.CreateUserDto true "User Data"
// @Success 201 {object} dtos.UserDto
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateUserDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		Confirmed: false,
		Role:      "USER",
	}

	createdUser, err := h.service.CreateUser(r.Context(), *user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// UpdateUser godoc
// @Summary Обновить данные пользователя
// @Description Обновление данных пользователя по ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body dtos.CreateUserDto true "Updated User Data"
// @Success 200 {object} dtos.UserDto
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Not found"
// @Failure 500 {object} string "Server error"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID пользователя из URL-параметра
	id := chi.URLParam(r, "id")

	// Декодируем тело запроса в структуру CreateUserDto
	var req dtos.CreateUserDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Создаем объект User из данных, полученных в запросе
	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	// Вызываем метод сервиса для обновления пользователя по ID
	updatedUser, err := h.service.UpdateUser(r.Context(), id, *user)
	if err != nil {
		// TODO: errors typization
		if err.Error() == "User not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем обновленного пользователя в ответе
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser godoc
// @Summary Удалить пользователя по ID
// @Description Удаление пользователя по ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 204 {string} string "User deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Not found"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Expected id url param", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteUser(ctx, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

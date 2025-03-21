package handlers

import (
	"encoding/json"
	"gateway/internal/dtos"
	"gateway/internal/models"
	"gateway/internal/services"
	"net/http"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Reftesh godoc
// @Summary Вход в систему
// @Description Авторизует пользователя по email и password. Выдает access и refresh токены.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param loginCredentials body dtos.LoginDto true "Данные для авторизации пользователя"
// @Success 201 {object} dtos.AuthCredentialsDto
// @Failure 400 {string} string "Неверные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginDto dtos.LoginDto
	if err := json.NewDecoder(r.Body).Decode(&loginDto); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	authCredentials, err := h.service.Login(r.Context(), &models.LoginCredentials{
		Email:    loginDto.Email,
		Password: loginDto.Password,
	})
	if err != nil {
		http.Error(w, "Ошибка при работе сервиса", http.StatusInternalServerError)
		return
	}

	authDto := dtos.AuthCredentialsDto{
		AccessToken:  authCredentials.AccessToken,
		RefreshToken: authCredentials.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authDto)
}

// Reftesh godoc
// @Summary Обновить сессионный токен
// @Description Обновляет сессионный access токен, инвалидирует старый refresh токен, создает новый.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param authCredentials body dtos.AuthCredentialsDto true "Данные сессионных токенов"
// @Success 201 {object} dtos.AuthCredentialsDto
// @Failure 400 {string} string "Неверные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var tokens dtos.AuthCredentialsDto
	if err := json.NewDecoder(r.Body).Decode(&tokens); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}
	authCredentials, err := h.service.Refresh(r.Context(), &models.AuthCredentials{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		http.Error(w, "Ошибка при работе сервиса", http.StatusInternalServerError)
		return
	}

	authDto := dtos.AuthCredentialsDto{
		AccessToken:  authCredentials.AccessToken,
		RefreshToken: authCredentials.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authDto)
}

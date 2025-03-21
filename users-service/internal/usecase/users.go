package usecase

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/utils"

	"github.com/google/uuid"
)

// UserService - сервис для работы с пользователями
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService - конструктор для создания сервиса пользователей
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create - создание нового пользователя
func (s *UserService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	// Генерация уникального ID для пользователя
	user.ID = uuid.NewString()
	var err error

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, user)
}

// GetByID - получение пользователя по ID
func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByID - получение пользователя по ID
func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// List - получение списка всех пользователей
func (s *UserService) List(ctx context.Context, page, limit int) ([]models.User, error) {
	return s.repo.List(ctx, page, limit)
}

// Update - обновление данных пользователя
func (s *UserService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	// Проверка, существует ли пользователь с таким ID
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Обновление данных пользователя
	existingUser.Email = user.Email
	existingUser.Username = user.Username
	existingUser.Password = user.Password
	existingUser.Confirmed = user.Confirmed
	existingUser.Role = user.Role
	existingUser.IsBlocked = user.IsBlocked

	// Сохранение обновленного пользователя в базе
	return s.repo.Update(ctx, existingUser)
}

// Delete - удаление пользователя по ID
func (s *UserService) Delete(ctx context.Context, id string) error {
	// Проверка, существует ли пользователь с таким ID
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	// Удаление пользователя из базы
	return s.repo.Delete(ctx, user.ID)
}

// Login - аутентификация пользователя
func (s *UserService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("пользователь не найден")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("неверный пароль")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	s.repo.SaveRefreshToken(ctx, &models.RefreshToken{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	})

	return accessToken, refreshToken, nil
}

// RefreshToken - обновление access-токена
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	userID, err := utils.ValidateToken(refreshToken, true)
	if err != nil {
		return "", "", errors.New("невалидный refresh-токен")
	}
	storedRefreshToken, err := s.repo.GetRefreshTokenByUserId(ctx, userID)
	if err != nil {
		return "", "", err
	}
	if storedRefreshToken != refreshToken {
		return "", "", errors.New("Invalid refresh token")
	}
	accessToken, refreshToken, err := utils.GenerateTokens(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken - проверка access-токена
func (s *UserService) ValidateToken(ctx context.Context, token string) (string, bool) {
	userID, err := utils.ValidateToken(token, false)
	if err != nil {
		return "", false
	}
	return userID, true
}

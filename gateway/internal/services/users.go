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

// UsersService - gRPC клиент для работы с пользователями
type UsersService struct {
	client proto.UserServiceClient
	logger *zap.Logger
}

// NewUsersService - конструктор для создания gRPC клиента
func NewUsersService(grpcAddress string, logger *zap.Logger) (*UsersService, error) {
	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Ошибка подключения к gRPC серверу", zap.String("address", grpcAddress), zap.Error(err))
		return nil, fmt.Errorf("не удалось подключиться к gRPC серверу: %w", err)
	}

	client := proto.NewUserServiceClient(conn)
	logger.Info("gRPC клиент успешно подключен", zap.String("address", grpcAddress))

	return &UsersService{client: client, logger: logger}, nil
}

// GetUsers - получение списка пользователей
func (s *UsersService) GetUsers(ctx context.Context, page, limit int) ([]models.User, error) {
	s.logger.Info("Запрос списка пользователей")

	resp, err := s.client.GetUsers(ctx, &proto.GetUsersRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		s.logger.Error("Ошибка получения списка пользователей", zap.Error(err))
		return nil, err
	}

	var users []models.User
	for _, u := range resp.Users {
		users = append(users, models.User{
			ID:        u.Id,
			Email:     u.Email,
			Username:  u.ProfileName,
			Password:  u.Password,
			Confirmed: u.Confirmed,
			Role:      u.Role,
			IsBlocked: u.IsBlocked,
		})
	}

	s.logger.Info("Список пользователей успешно получен", zap.Int("count", len(users)))
	return users, nil
}

// GetUserByID - получение пользователя по ID
func (s *UsersService) GetUserByID(ctx context.Context, id string) (models.User, error) {
	s.logger.Info("Запрос пользователя по ID", zap.String("id", id))

	resp, err := s.client.GetUserByID(ctx, &proto.GetUserByIDRequest{Id: id})
	if err != nil {
		s.logger.Error("Ошибка получения пользователя", zap.String("id", id), zap.Error(err))
		return models.User{}, err
	}

	user := models.User{
		ID:        resp.User.Id,
		Email:     resp.User.Email,
		Username:  resp.User.ProfileName,
		Password:  resp.User.Password,
		Confirmed: resp.User.Confirmed,
		Role:      resp.User.Role,
		IsBlocked: resp.User.IsBlocked,
	}

	s.logger.Info("Пользователь успешно получен", zap.String("id", user.ID))
	return user, nil
}

// CreateUser - создание нового пользователя
func (s *UsersService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	s.logger.Info("Создание нового пользователя", zap.String("email", user.Email))

	resp, err := s.client.CreateUser(ctx, &proto.CreateUserRequest{
		Email:       user.Email,
		ProfileName: user.Username,
		Password:    user.Password,
		Confirmed:   user.Confirmed,
		Role:        user.Role,
	})
	if err != nil {
		s.logger.Error("Ошибка создания пользователя", zap.String("email", user.Email), zap.Error(err))
		return models.User{}, err
	}

	createdUser := models.User{
		ID:        resp.Id,
		Email:     resp.Email,
		Username:  resp.ProfileName,
		Password:  resp.Password,
		Confirmed: resp.Confirmed,
		Role:      resp.Role,
		IsBlocked: resp.IsBlocked,
	}

	s.logger.Info("Пользователь успешно создан", zap.String("id", createdUser.ID))
	return createdUser, nil
}

// UpdateUser - обновление пользователя
func (s *UsersService) UpdateUser(ctx context.Context, id string, user models.User) (models.User, error) {
	s.logger.Info("Обновление пользователя", zap.String("id", id))

	resp, err := s.client.UpdateUser(ctx, &proto.UpdateUserRequest{
		Id:          id,
		ProfileName: user.Username,
		Password:    user.Password,
	})
	if err != nil {
		s.logger.Error("Ошибка обновления пользователя", zap.String("id", id), zap.Error(err))
		return models.User{}, err
	}

	updatedUser := models.User{
		ID:        resp.Id,
		Email:     resp.Email,
		Username:  resp.ProfileName,
		Password:  resp.Password,
		Confirmed: resp.Confirmed,
		Role:      resp.Role,
		IsBlocked: resp.IsBlocked,
	}

	s.logger.Info("Пользователь успешно обновлен", zap.String("id", updatedUser.ID))
	return updatedUser, nil
}

// DeleteUser - удаление пользователя
func (s *UsersService) DeleteUser(ctx context.Context, id string) error {
	s.logger.Info("Удаление пользователя", zap.String("id", id))

	resp, err := s.client.DeleteUser(ctx, &proto.DeleteUserRequest{Id: id})
	if err != nil {
		s.logger.Error("Ошибка удаления пользователя", zap.String("id", id), zap.Error(err))
		return err
	}

	if !resp.Success {
		s.logger.Warn("Пользователь не был удален", zap.String("id", id))
		return fmt.Errorf("не удалось удалить пользователя")
	}

	s.logger.Info("Пользователь успешно удален", zap.String("id", id))
	return nil
}

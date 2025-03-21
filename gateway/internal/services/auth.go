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

type AuthService struct {
	client proto.UserServiceClient
	logger *zap.Logger
}

// NewAuthService - конструктор сервиса с логированием
func NewAuthService(grpcAddress string, logger *zap.Logger) (*AuthService, error) {
	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Ошибка подключения к gRPC серверу", zap.String("address", grpcAddress), zap.Error(err))
		return nil, fmt.Errorf("не удалось подключиться к gRPC серверу: %w", err)
	}

	client := proto.NewUserServiceClient(conn)
	logger.Info("gRPC клиент успешно подключен", zap.String("address", grpcAddress))

	return &AuthService{client: client, logger: logger}, nil
}

func (s *AuthService) Login(ctx context.Context, credentials *models.LoginCredentials) (*models.AuthCredentials, error) {
	s.logger.Info("Авторизация пользователя", zap.String("email", credentials.Email))

	resp, err := s.client.Login(ctx, &proto.LoginRequest{
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		s.logger.Error("Ошибка авторизации", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Авторизация пользователя прошла успешно", zap.String("email", credentials.Email))
	return &models.AuthCredentials{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, authCredentials *models.AuthCredentials) (*models.AuthCredentials, error) {
	s.logger.Info("Обновление сессионного токена")

	resp, err := s.client.RefreshToken(ctx, &proto.RefreshTokenRequest{
		RefreshToken: authCredentials.RefreshToken,
	})
	if err != nil {
		s.logger.Error("Ошибка обновления сессионного токена", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Обновление сессионного токена прошло успешно")
	return &models.AuthCredentials{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
	}, nil
}

func (s *AuthService) Validate(ctx context.Context, access string) (bool, error) {
	s.logger.Info("Валидация сессионного токена")

	resp, err := s.client.ValidateToken(ctx, &proto.ValidateTokenRequest{
		AccessToken: access,
	})
	if err != nil {
		s.logger.Error("Ошибка валидации сессионного токена", zap.Error(err))
		return false, err
	}
	s.logger.Info("Валидация сессионного токена прошла успешно")
	return resp.GetValid(), nil
}

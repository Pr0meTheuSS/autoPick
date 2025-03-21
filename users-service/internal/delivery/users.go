package delivery

import (
	"context"
	"user-service/internal/models"
	"user-service/internal/proto"
	"user-service/internal/usecase"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ proto.UserServiceServer = (*UserHandler)(nil)

// UserHandler - обработчик запросов для пользователей
type UserHandler struct {
	proto.UnimplementedUserServiceServer
	service *usecase.UserService
	logger  *zap.Logger
}

// NewUserHandler - конструктор для создания обработчика
func NewUserHandler(service *usecase.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

// CreateUser - обработка запроса на создание пользователя
func (h *UserHandler) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
	h.logger.Info("Получен запрос на создание пользователя", zap.String("email", req.Email))

	user := &models.User{
		Email:     req.Email,
		Username:  req.ProfileName,
		Password:  req.Password,
		Confirmed: req.Confirmed,
		Role:      req.Role,
	}

	createdUser, err := h.service.Create(ctx, user)
	if err != nil {
		h.logger.Error("Ошибка при создании пользователя", zap.String("email", req.Email), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось создать пользователя: %v", err)
	}

	h.logger.Info("Пользователь успешно создан", zap.String("id", createdUser.ID))

	return &proto.User{
		Id:          createdUser.ID,
		Email:       createdUser.Email,
		ProfileName: createdUser.Username,
		Password:    createdUser.Password,
		Confirmed:   createdUser.Confirmed,
		Role:        createdUser.Role,
		IsBlocked:   createdUser.IsBlocked,
	}, nil
}

// GetUserByID - обработка запроса на получение пользователя по ID
func (h *UserHandler) GetUserByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.GetUserByIDResponse, error) {
	h.logger.Info("Получен запрос на получение пользователя", zap.String("id", req.Id))

	user, err := h.service.GetByID(ctx, req.Id)
	if err != nil {
		h.logger.Error("Пользователь не найден", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "не удалось получить пользователя с ID %s: %v", req.Id, err)
	}

	h.logger.Info("Пользователь успешно найден", zap.String("id", user.ID))

	return &proto.GetUserByIDResponse{
		User: &proto.User{
			Id:          user.ID,
			Email:       user.Email,
			ProfileName: user.Username,
			Password:    user.Password,
			Confirmed:   user.Confirmed,
			Role:        user.Role,
			IsBlocked:   user.IsBlocked,
		},
	}, nil
}

// GetUsers - обработка запроса на получение списка пользователей
func (h *UserHandler) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	h.logger.Info("Получен запрос на получение списка пользователей", zap.Int("page", int(req.Page)), zap.Int("limit", int(req.Limit)))

	users, err := h.service.List(ctx, int(req.Page), int(req.Limit))
	if err != nil {
		h.logger.Error("Ошибка при получении списка пользователей", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось получить список пользователей: %v", err)
	}

	h.logger.Info("Список пользователей успешно получен", zap.Int("count", len(users)))

	var userList []*proto.User
	for _, user := range users {
		userList = append(userList, &proto.User{
			Id:          user.ID,
			Email:       user.Email,
			ProfileName: user.Username,
			Password:    user.Password,
			Confirmed:   user.Confirmed,
			Role:        user.Role,
			IsBlocked:   user.IsBlocked,
		})
	}

	return &proto.GetUsersResponse{
		Users: userList,
	}, nil
}

// UpdateUser - обработка запроса на обновление пользователя
func (h *UserHandler) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.User, error) {
	h.logger.Info("Получен запрос на обновление пользователя", zap.String("id", req.Id))

	user := &models.User{
		ID:       req.Id,
		Username: req.ProfileName,
		Password: req.Password,
	}

	updatedUser, err := h.service.Update(ctx, user)
	if err != nil {
		h.logger.Error("Ошибка при обновлении пользователя", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось обновить пользователя: %v", err)
	}

	h.logger.Info("Пользователь успешно обновлен", zap.String("id", req.Id))

	return &proto.User{
		Id:          updatedUser.ID,
		Email:       updatedUser.Email,
		ProfileName: updatedUser.Username,
		Password:    updatedUser.Password,
		Confirmed:   updatedUser.Confirmed,
		Role:        updatedUser.Role,
		IsBlocked:   updatedUser.IsBlocked,
	}, nil
}

// DeleteUser - обработка запроса на удаление пользователя
func (h *UserHandler) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	h.logger.Info("Получен запрос на удаление пользователя", zap.String("id", req.Id))

	err := h.service.Delete(ctx, req.Id)
	if err != nil {
		h.logger.Error("Ошибка при удалении пользователя", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "не удалось удалить пользователя с ID %s: %v", req.Id, err)
	}

	h.logger.Info("Пользователь успешно удален", zap.String("id", req.Id))

	return &proto.DeleteUserResponse{
		Success: true,
	}, nil
}

// Login - обработка входа
func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	accessToken, refreshToken, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Warn("Ошибка аутентификации", zap.String("email", req.Email), zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "неверные учетные данные")
	}

	return &proto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken - обновление access-токена
func (h *UserHandler) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := h.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "невалидный refresh-токен")
	}

	return &proto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken - проверка access-токена
func (h *UserHandler) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	userID, valid := h.service.ValidateToken(ctx, req.AccessToken)
	return &proto.ValidateTokenResponse{Valid: valid, UserId: userID}, nil
}

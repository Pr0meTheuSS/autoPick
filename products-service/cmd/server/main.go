package main

import (
	"context"
	"log"
	"net"
	"product-service/internal/delivery"
	"product-service/internal/proto" // Путь к вашему сгенерированному файлу
	"product-service/internal/repository"
	"product-service/internal/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap" // Импортируем zap для логирования
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Создаем логгер
	logger, err := zap.NewProduction() // Создаем стандартный логгер
	if err != nil {
		log.Fatalf("Ошибка при создании логгера: %v", err)
	}
	defer logger.Sync() // Закрытие логгера при завершении работы программы

	// Подключаемся к MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.Fatal("Ошибка при подключении к MongoDB", zap.Error(err))
	}

	// Проверяем соединение с базой данных
	err = client.Ping(context.Background(), nil)
	if err != nil {
		logger.Fatal("Ошибка при проверке соединения с MongoDB", zap.Error(err))
	}

	// Получаем доступ к нужной базе данных
	db := client.Database("productDB") // Используйте имя вашей базы данных

	// Создаем gRPC сервер
	server := grpc.NewServer()

	// Создаем репозиторий, сервис и обработчик
	repository := repository.NewProductRepository(db)
	service := usecase.NewProductService(repository)
	handler := delivery.NewProductHandler(service, logger) // Передаем логгер в обработчик

	// Регистрируем сервис (например, ProductService)
	proto.RegisterProductServiceServer(server, handler)

	// Включаем рефлексию
	reflection.Register(server)

	// Настроим и запустим сервер
	listener, err := net.Listen("tcp", ":9091")
	if err != nil {
		logger.Fatal("Ошибка при создании слушателя", zap.Error(err))
	}

	logger.Info("Сервер запущен на :9091")
	if err := server.Serve(listener); err != nil {
		logger.Fatal("Ошибка при запуске сервера", zap.Error(err))
	}
}

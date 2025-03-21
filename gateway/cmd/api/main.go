package main

import (
	"fmt"
	"gateway/docs"
	"gateway/internal/handlers"
	middlewares "gateway/internal/middleware"
	"gateway/internal/services"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Gateway Service API
// @version 1.0
// @description gateway API
// @host localhost:9090
// @BasePath /
func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить .env, используем переменные окружения")
	}
	// Подключаемся к gRPC-серверу
	// TODO: get host:port from envs
	// Создаем логгер Zap (в проде лучше использовать zap.NewProduction())
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer logger.Sync()

	// Создаем сервис с логированием
	grpcAddress := "localhost:9091"
	productService, err := services.NewProductsService(grpcAddress, logger)
	if err != nil {
		logger.Fatal("Ошибка создания gRPC клиента", zap.Error(err))
	}

	logger.Info("Приложение запущено и готово принимать запросы")

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cookie"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           500,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Получаем текущий `basePath` из переменной окружения
	basePath := os.Getenv("BASE_URL")
	if basePath == "" {
		basePath = "localhost:9090" // Значение по умолчанию
	}
	docs.SwaggerInfo.Host = basePath
	productHandler := handlers.NewProductsHandler(*productService)

	r.Route("/products", func(r chi.Router) {
		r.With(middlewares.PaginationMiddleware).Get("/", productHandler.Get)
		r.Get("/{id}", productHandler.GetByID)
		// TODO: add protect middleware
		r.Post("/", productHandler.Post)
		// TODO: add protect middleware
		r.Delete("/{id}", productHandler.Delete)
		// TODO: add protect middleware
		r.Put("/{id}", productHandler.Put)
	})

	// TODO: handle error
	userService, _ := services.NewUsersService("localhost:9092", logger)
	userHandler := handlers.NewUserHandler(userService)
	r.Route("/users", func(r chi.Router) {
		// TODO: only admin middleware
		r.With(middlewares.PaginationMiddleware).Get("/", userHandler.GetUsers) // Получить всех пользователей
		r.Get("/{id}", userHandler.GetUserByID)                                 // Получить пользователя по ID
		r.Post("/", userHandler.CreateUser)                                     // Создать нового пользователя
		r.Put("/{id}", userHandler.UpdateUser)                                  // Обновить данные пользователя
		// TODO: only admin middleware
		r.Delete("/{id}", userHandler.DeleteUser) // Удалить пользователя
		// TODO: Add block, confirm handlers
	})
	authService, _ := services.NewAuthService("localhost:9092", logger)
	authHandler := handlers.NewAuthHandler(authService)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
	})

	orderService, _ := services.NewOrdersService("localhost:9093", logger)
	orderHandler := handlers.NewOrdersHandler(*orderService)

	r.Route("/orders", func(r chi.Router) {
		// TODO: add pagination middleware
		r.With(middlewares.PaginationMiddleware).Get("/", orderHandler.Get)

		r.Get("/{id}", orderHandler.GetByID)
		// TODO: add protect middleware
		r.Post("/", orderHandler.Post)
		// TODO: add protect middleware
		r.Delete("/{id}", orderHandler.Delete)
		// TODO: add protect middleware
		r.Put("/{id}", orderHandler.Put)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Запуск сервера
	port := getEnv("CONTENT_SERVICE_PORT", "9090")

	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Функция для получения переменной окружения с дефолтным значением
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

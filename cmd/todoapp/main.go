// Точка входа приложения. Здесь происходит:
//   - Инициализация конфигурации и логгера
//   - Подключение к базе данных PostgreSQL
//   - «Сборка» всех фич (Repository → Service → HTTP Handler)
//   - Запуск HTTP-сервера с graceful shutdown
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/ArthasEden/todo-app/internal/core/logger"
	core_postgres_pool "github.com/ArthasEden/todo-app/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/ArthasEden/todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/ArthasEden/todo-app/internal/core/transport/http/server"
	users_postgres_repository "github.com/ArthasEden/todo-app/internal/features/users/repository/postgres"
	users_service "github.com/ArthasEden/todo-app/internal/features/users/service"
	users_transport_http "github.com/ArthasEden/todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	// Создаём корневой контекст, который отменяется при получении SIGINT/SIGTERM
	// (Ctrl+C или команда `kill`). Это основа для graceful shutdown.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Инициализируем логгер приложения
	// Пишет одновременно в stdout и в файл (см. internal/core/logger).
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	// Создаём пулл соединений с PostgreSQL через библиотеку pgx.
	// Пул переиспользует соединения, что гораздо эффективнее,
	// чем открывать новое соединение на каждый SQL запрос.
	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	// Ручное внедрение зависимостей (Dependency Injection):
	// Repository → Service → HTTP Handler.
	// Каждый слой знает только об интерфейсе нижележащего.
	// Это обеспечивает слабую связанность (loose coupling) и тестируемость.
	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	// Собираем HTTP-сервер с цепочкой middleware.
	// Middleware применяются ко всем маршрутам (Route) в порядке объявления:
	// CORS → RequestID → Logger → Trace → Panic recovery.
	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	// Регистрируем маршруты API v1.
	// APIVersionRouter автоматически добавляет префикс /api/v1 ко всем путям.
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	httpServer.RegisterAPIRoutes(*apiVersionRouter)

	// Запускаем сервер. Блокируется до получения сигнала завершения.
	// После сигнала выполняет graceful shutdown: ждёт завершения активных запросов.
	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}

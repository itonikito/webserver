package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webserver/config"
	"webserver/internal/handlers"
	"webserver/internal/middleware"
	"webserver/internal/repositories"
	"webserver/internal/services"
	"webserver/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация логгера
	if err := logger.Init(cfg); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	log := logger.Get()
	log.Info("starting application",
		logger.String("port", cfg.ServerPort),
		logger.String("version", "1.0.0"))

	// Инициализация зависимостей
	batRepo := repositories.NewBatRepository(cfg.BatFilePath)
	batService := services.NewBatService(batRepo, &services.ServiceConfig{
		TimeoutSec:    cfg.TimeoutSec,
		MaxConcurrent: int64(cfg.MaxConcurrent),
	})
	batHandler := handlers.NewBatHandler(batService)

	// Настройка маршрутизатора
	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryMiddleware)

	// Регистрация маршрутов
	batHandler.RegisterRoutes(router)

	// Настройка HTTP сервера
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		log.Info("server is starting", logger.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", logger.Error(err))
		}
	}()

	// Ожидание сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", logger.Error(err))
	}

	log.Info("server stopped")
}

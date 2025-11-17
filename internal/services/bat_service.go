package services

import (
	"context"
	"fmt"
	"time"
	"webserver/internal/entities"
	"webserver/internal/repositories"
	"webserver/pkg/errors"

	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type BatService interface {
	ProcessRequest(ctx context.Context, action string, param string) (*entities.BatResponse, error)
	ValidateAction(action string) bool
	GetStats() *ServiceStats
}

type batService struct {
	batRepo   repositories.BatRepository
	timeout   time.Duration
	semaphore *semaphore.Weighted
	config    *ServiceConfig
	stats     *ServiceStats
}

type ServiceConfig struct {
	TimeoutSec    int
	MaxConcurrent int64
}

type ServiceStats struct {
	TotalRequests   int64 `json:"total_requests"`
	SuccessRequests int64 `json:"success_requests"`
	FailedRequests  int64 `json:"failed_requests"`
	ActiveRequests  int64 `json:"active_requests"`
}

func NewBatService(batRepo repositories.BatRepository, config *ServiceConfig) BatService {
	return &batService{
		batRepo:   batRepo,
		timeout:   time.Duration(config.TimeoutSec) * time.Second,
		semaphore: semaphore.NewWeighted(config.MaxConcurrent),
		config:    config,
		stats:     &ServiceStats{},
	}
}

func (s *batService) ProcessRequest(ctx context.Context, action string, param string) (*entities.BatResponse, error) {
	// Обновляем статистику
	s.stats.TotalRequests++
	s.stats.ActiveRequests++
	defer func() { s.stats.ActiveRequests-- }()

	// Ограничиваем concurrent запросы
	if err := s.semaphore.Acquire(ctx, 1); err != nil {
		s.stats.FailedRequests++
		return nil, errors.NewTimeoutError("too many concurrent requests")
	}
	defer s.semaphore.Release(1)

	// Валидация действия
	if !s.ValidateAction(action) {
		s.stats.FailedRequests++
		return nil, errors.NewValidationError(map[string]interface{}{
			"action":  action,
			"allowed": []string{"create", "update"},
		})
	}

	// Устанавливаем таймаут
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Логируем начало выполнения
	zap.L().Info("starting bat execution",
		zap.String("action", action),
		zap.String("param", param))

	// Выполняем .bat файл
	startTime := time.Now()
	output, err := s.batRepo.ExecuteBat(ctx, action, param)
	executionTime := time.Since(startTime)

	if err != nil {
		s.stats.FailedRequests++
		zap.L().Error("bat execution failed",
			zap.String("action", action),
			zap.String("param", param),
			zap.Duration("duration", executionTime),
			zap.Error(err))
		return nil, err
	}

	s.stats.SuccessRequests++
	zap.L().Info("bat execution completed",
		zap.String("action", action),
		zap.String("param", param),
		zap.Duration("duration", executionTime))

	// Формируем успешный ответ
	message := getActionDescription(action, param)

	return &entities.BatResponse{
		Success:    true,
		Message:    message,
		Output:     output,
		Action:     action,
		Parameter:  param,
		DurationMs: executionTime.Milliseconds(),
		Timestamp:  time.Now(),
	}, nil
}

func (s *batService) ValidateAction(action string) bool {
	validActions := map[string]bool{
		"create": true,
		"update": true,
	}
	return validActions[action]
}

func (s *batService) GetStats() *ServiceStats {
	return s.stats
}

func getActionDescription(action, param string) string {
	switch action {
	case "create":
		return fmt.Sprintf("База данных '%s' успешно создана", param)
	case "update":
		return fmt.Sprintf("База данных '%s' успешно обновлена", param)
	default:
		return fmt.Sprintf("Операция '%s' для базы '%s' выполнена", action, param)
	}
}

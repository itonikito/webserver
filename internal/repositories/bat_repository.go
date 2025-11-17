package repositories

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"webserver/pkg/errors"

	"go.uber.org/zap"
)

type BatRepository interface {
	ExecuteBat(ctx context.Context, action string, param string) (string, error)
}

type batRepository struct {
	batFilePath string
}

func NewBatRepository(batFilePath string) BatRepository {
	return &batRepository{
		batFilePath: batFilePath,
	}
}

func (r *batRepository) ExecuteBat(ctx context.Context, action string, param string) (string, error) {
	// Получаем абсолютный путь к файлу для лучшей диагностики
	absPath, err := filepath.Abs(r.batFilePath)
	if err != nil {
		absPath = r.batFilePath
	}

	// Проверяем существование файла
	if !r.fileExists(absPath) {
		zap.L().Error("batch file not found",
			zap.String("requested_path", r.batFilePath),
			zap.String("absolute_path", absPath),
			zap.String("working_dir", r.getWorkingDir()))

		return "", errors.NewBatExecutionError(
			map[string]interface{}{
				"requested_path": r.batFilePath,
				"absolute_path":  absPath,
				"working_dir":    r.getWorkingDir(),
			},
			fmt.Errorf("batch file not found"),
		)
	}

	zap.L().Debug("batch file found",
		zap.String("path", absPath),
		zap.String("action", action),
		zap.String("param", param))

	// Запускаем .bat файл с действием и параметром
	//cmd := exec.CommandContext(ctx, "cmd", "/C", absPath, action, param)
	// Для Mac/Linux используем bash
	cmd := exec.CommandContext(ctx, "/bin/bash", absPath, action, param)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Проверяем, не был ли отменен контекст (таймаут)
		if ctx.Err() == context.DeadlineExceeded {
			zap.L().Warn("bat execution timeout",
				zap.String("action", action),
				zap.String("param", param))
			return "", errors.NewTimeoutError(map[string]interface{}{
				"action": action,
				"param":  param,
			})
		}

		zap.L().Error("bat execution failed",
			zap.String("action", action),
			zap.String("param", param),
			zap.String("output", string(output)),
			zap.Error(err))

		return "", errors.NewBatExecutionError(
			map[string]interface{}{
				"action": action,
				"param":  param,
				"file":   absPath,
				"output": string(output),
			},
			err,
		)
	}

	zap.L().Debug("bat execution completed successfully",
		zap.String("action", action),
		zap.String("param", param))

	return string(output), nil
}

func (r *batRepository) fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (r *batRepository) getWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return wd
}

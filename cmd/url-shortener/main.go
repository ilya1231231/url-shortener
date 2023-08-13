package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"os"
	"url-shortender/internal/config"
	"url-shortender/internal/lib/sl"
	"url-shortender/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//данная библиотека минималистичная и может читать почти со всех файлов конфига (json, yaml, env)
	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env)
	logger.Info("starting url-shortener", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")
	//@todo сделать подключение к другой СУБД
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	err = storage.DeleteUrl("www")
	if err != nil {
		logger.Error("failed to delete url", sl.Err(err))
		os.Exit(1)
	}
	router := chi.NewRouter()
	//у chi есть мидлвары из коробки
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	//нужен для того, чтоы при панике в хендлере не падало все приложение
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	/*	для вывода в файл
		file, err := os.Create("log.txt")
		if err != nil {
			log.Fatal("Can`t create log file")
		}
	*/
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		//минимальный уровень выводимых ошибок ошибки до этой константы игнорируются
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}

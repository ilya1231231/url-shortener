package main

import (
	"golang.org/x/exp/slog"
	"log"
	"os"
	"url-shortender/internal/config"
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
	//вот тут я бы использовал mysql
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		//slog как я понял пока не умеет обрабатывать в таком режиме, нужно писать отдельнуюфункцию
		log.Fatalf("failed to init storage %s", err.Error())
	}
	_ = storage
	// TODO: init router: chi, "chi/render"

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

package main

import (
	"golang.org/x/exp/slog"
	"os"
	"url-shortender/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//данная библиотека минималистичная и может читать почти со всех файлов конфига (json, yaml, env)
	// TODO: init config: cleanenv
	//!!! изучть как правильно используется Setenv
	os.Setenv("CONFIG_PATH", "./config/local.yaml")
	cfg := config.MustLoad()
	//библиотека для логирования
	//изучить и понять логику детальнее. При env=prod не выводится сообщение из-за настройки.
	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	//вот тут я бы использовал mysql
	// TODO: init storage: sqlite

	//
	// TODO: init router: chi, "chi/render"

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		//минимальный уровень выводимых ошибок ошибки до этой константы игнорируются
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

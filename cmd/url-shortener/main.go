package main

import (
	"fmt"
	"os"
	"url-shortender/internal/config"
)

func main() {
	//данная библиотека минималистичная и может читать почти со всех файлов конфига (json, yaml, env)
	// TODO: init config: cleanenv
	os.Setenv("CONFIG_PATH", "./config/local.yaml")
	cfg := config.MustLoad()
	fmt.Println(cfg)
	//библиотека для логирования
	// TODO: init logger: slog

	//вот тут я бы использовал mysql
	// TODO: init storage: sqlite

	//
	// TODO: init router: chi, "chi/render"

	// TODO: run server
}

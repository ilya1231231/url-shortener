package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	//struct tag
	//структурные теги - это аннотации, показывающие какую переменную присваивать полю структуры при парсинге, например yaml файла
	//также можно назначить переменную с варсинга json (`json:env`)
	//Проще- сопоставление параметров в объекте и параметров в yaml файле
	Env         string `yaml:"env" env-default:"dev"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8082"`
	Timout      time.Duration `yaml:"timout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// функции с префикса Must применяется при панике
func MustLoad() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH not found")
	}

	//os.Stat возвращает FileInfo
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file by path name %s does not exists", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
}

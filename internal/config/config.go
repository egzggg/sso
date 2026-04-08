package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath() // сохраняем в переменную путь до файла

	if path == "" {
		panic("config file does not exist: " + path)
	}

	return MustLoadByPath(path) // из нашего ямл файла возвращаем структуру с заполненными полями
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file") //записываем значение из флага в переменную, тоесть путь в нашемслучпе
	flag.Parse()  // бьет команду строки, ищет флаги и записывает значение

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res // итоговый путь
}

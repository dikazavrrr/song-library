package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config содержит конфигурационные параметры приложения
type Config struct {
	Postgres   Postgres `mapstructure:"postgres"`
	HTTPServer HTTPConfig     `mapstructure:"httpserver"`
}

// Postgres хранит настройки подключения к БД
type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// HTTPConfig хранит настройки HTTP-сервера
type HTTPConfig struct {
	Address string `mapstructure:"address"`
}

// MustInit загружает конфигурацию
func MustInit(isProd string) Config {
	viper.SetConfigName("config")   // Имя файла (без расширения)
	viper.SetConfigType("yaml")     // Тип файла
	viper.AddConfigPath("./config") // Папка с конфигом

	// Читаем конфиг
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка чтения конфига: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Ошибка декодирования конфига: %v", err)
	}

	return cfg
}

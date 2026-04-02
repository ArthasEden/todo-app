package core_logger // тот же пакет конфигурации логгера

import (
	"fmt" // форматирование ошибок

	"github.com/kelseyhightower/envconfig" // библиотека для чтения env в struct
)

// Config — структура конфигурации логгера
type Config struct {
	Level  string `envconfig:"LEVEL"  required:"true"` // уровень логов (например INFO)
	Folder string `envconfig:"FOLDER" required:"true"` // папка для логов
}

// NewConfig — загружает конфиг из env
func NewConfig() (Config, error) {
	var config Config // создаём пустую структуру

	// читаем переменные окружения с префиксом LOGGER_
	// например LOGGER_LEVEL, LOGGER_FOLDER
	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil // возвращаем заполненный конфиг
}

// NewConfigMust — вариант, который падает при ошибке
func NewConfigMust() Config {
	config, err := NewConfig() // пробуем загрузить конфиг
	if err != nil {
		err = fmt.Errorf("get Logger config: %w", err) // добавляем контекст
		panic(err)                                     // останавливаем приложение
	}

	return config // возвращаем конфиг
}

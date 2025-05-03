package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rx3lixir/agg-api/internal/lib/logger"
	"github.com/spf13/viper"
)

// Константы для ключей конфигурации
const (
	envKey            = "application_params.env"
	usernameKey       = "db_params.username"
	passwordKey       = "db_params.password"
	dbNameKey         = "db_params.db_name"
	hostKey           = "db_params.host"
	portKey           = "db_params.port"
	connectTimeoutKey = "db_params.connect_timeout"
	secretKey         = "server_params.secret_key"
)

// AppConfig представляет конфигурацию всего приложения
type AppConfig struct {
	Application ApplicationParams `mapstructure:"application_params" validate:"required"`
	DB          DBParams          `mapstructure:"db_params" validate:"required"`
	Server      ServerParams      `mapstructure:"server_params" validate:"required"`
}

// ApplicationParams содержит общие параметры приложения
type ApplicationParams struct {
	Env string `mapstructure:"env" validate:"required,oneof=dev prod test"`
}

type ServerParams struct {
	Address   string `mapstructure:"address" validate:"required"`
	SecretKey string `mapstructure:"secret_key" validate:"required"`
}

// DBParams содержит параметры подключения к базе данных
type DBParams struct {
	Username       string        `mapstructure:"username" validate:"required"`
	Password       string        `mapstructure:"password" validate:"required"`
	DBName         string        `mapstructure:"db_name" validate:"required"`
	Host           string        `mapstructure:"host" validate:"required"`
	Port           int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	ConnectTimeout time.Duration `mapstructure:"connect_timeout" validate:"required,min=1"`
}

// DSN собирает строку подключения к базе данных
func (db *DBParams) DSN() string {
	// Если хост не указан, используем localhost по умолчанию
	host := db.Host
	if host == "" {
		host = "localhost"
	}

	// Преобразование timeout в секунды
	timeoutSec := int(db.ConnectTimeout.Seconds())
	if timeoutSec < 1 {
		timeoutSec = 10 // Значение по умолчанию, если timeout некорректный
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?connect_timeout=%d&sslmode=disable",
		db.Username,
		db.Password,
		host,
		db.Port,
		db.DBName,
		timeoutSec,
	)
}

// EnvBindings возвращает мапу ключей конфигурации и соответствующих им переменных окружения
func envBindings() map[string]string {
	return map[string]string{
		envKey:            "APP_ENV",
		usernameKey:       "DB_USERNAME",
		passwordKey:       "DB_PASSWORD",
		dbNameKey:         "DB_NAME",
		hostKey:           "DB_HOST",
		portKey:           "DB_PORT",
		connectTimeoutKey: "DB_CONNECT_TIMEOUT",
		secretKey:         "SECRET_KEY",
	}
}

// New загружает конфигурацию из файла и переменных окружения
func New(log logger.Logger) (*AppConfig, error) {
	v := viper.New()

	// Настройка путей и формата
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	// Привязка переменных окружения
	for configKey, envVar := range envBindings() {
		if err := v.BindEnv(configKey, envVar); err != nil {
			return nil, fmt.Errorf("ошибка привязки переменной окружения %s: %w", envVar, err)
		}
	}

	// Чтение конфигурации
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурационного файла: %w", err)
	}

	var config AppConfig

	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании конфигурации: %w", err)
	}

	// Установка значений по умолчанию
	if config.DB.Host == "" {
		config.DB.Host = "localhost"
	}

	// Валидация конфигурации
	validate := validator.New()

	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	log.Info("Конфигурация успешно загружена")
	return &config, nil
}

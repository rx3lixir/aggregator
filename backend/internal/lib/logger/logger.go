package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

// LogLevel представляет уровень логирования
type LogLevel string

// LogFormat представляет формат вывода логов
type LogFormat string

// Константы для уровней логирования
const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)

// Константы для форматов логирования
const (
	TextFormat LogFormat = "text"
	JSONFormat LogFormat = "json"
)

// Logger определяет интерфейс для логирования
// Это позволяет отвязаться от конкретной библиотеки
type Logger interface {
	Debug(msg interface{}, keyvals ...interface{})
	Info(msg interface{}, keyvals ...interface{})
	Warn(msg interface{}, keyvals ...interface{})
	Error(msg interface{}, keyvals ...interface{})
	Fatal(msg interface{}, keyvals ...interface{})
	With(keyvals ...interface{}) Logger // Позволяет создавать дочерние логгеры с доп. полями
}

// charmLogger реализует интерфейс Logger с использованием charmbracelet/log
type charmLogger struct {
	logger *log.Logger
}

// Проверка на этапе компиляции, что тип реализует интерфейс
var _ Logger = (*charmLogger)(nil)

// Debug логирует сообщение на уровне Debug
func (l *charmLogger) Debug(msg interface{}, keyvals ...interface{}) {
	l.logger.Debug(msg, keyvals...)
}

// Info логирует сообщение на уровне Info
func (l *charmLogger) Info(msg interface{}, keyvals ...interface{}) {
	l.logger.Info(msg, keyvals...)
}

// Warn логирует сообщение на уровне Warn
func (l *charmLogger) Warn(msg interface{}, keyvals ...interface{}) {
	l.logger.Warn(msg, keyvals...)
}

// Error логирует сообщение на уровне Error
func (l *charmLogger) Error(msg interface{}, keyvals ...interface{}) {
	l.logger.Error(msg, keyvals...)
}

// Fatal логирует сообщение на уровне Fatal
func (l *charmLogger) Fatal(msg interface{}, keyvals ...interface{}) {
	l.logger.Fatal(msg, keyvals...)
}

// With создает новый логгер с дополнительным контекстом
func (l *charmLogger) With(keyvals ...interface{}) Logger {
	return &charmLogger{
		logger: l.logger.With(keyvals...),
	}
}

// Config содержит настройки для логгера
type Config struct {
	Level  LogLevel  // Уровень логирования
	Format LogFormat // Формат вывода
	Output io.Writer // Куда писать логи
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		Level:  InfoLevel,
		Format: TextFormat,
		Output: os.Stderr,
	}
}

// New создает новый логгер на основе предоставленной конфигурации
func New(cfg Config) Logger {
	// Применяем значения по умолчанию, если не заданы
	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}

	// Определяем уровень логирования
	var level log.Level
	switch strings.ToLower(string(cfg.Level)) {
	case string(DebugLevel):
		level = log.DebugLevel
	case string(InfoLevel):
		level = log.InfoLevel
	case string(WarnLevel):
		level = log.WarnLevel
	case string(ErrorLevel):
		level = log.ErrorLevel
	default:
		level = log.InfoLevel
	}

	// Определяем формат логирования
	var format log.Formatter
	switch strings.ToLower(string(cfg.Format)) {
	case string(JSONFormat):
		format = log.JSONFormatter
	case string(TextFormat):
		format = log.TextFormatter
	default:
		format = log.TextFormatter
	}

	// Создаем логгер
	logger := log.NewWithOptions(cfg.Output, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.RFC3339,
		Level:           level,
		Formatter:       format,
	})

	return &charmLogger{
		logger: logger,
	}
}

// NewDebug создает новый логгер с уровнем Debug
func NewDebug(output io.Writer) Logger {
	cfg := DefaultConfig()
	cfg.Level = DebugLevel
	if output != nil {
		cfg.Output = output
	}
	return New(cfg)
}

// NewJSON создает новый логгер с JSON форматом
func NewJSON(level LogLevel, output io.Writer) Logger {
	cfg := DefaultConfig()
	cfg.Format = JSONFormat
	cfg.Level = level
	if output != nil {
		cfg.Output = output
	}
	return New(cfg)
}

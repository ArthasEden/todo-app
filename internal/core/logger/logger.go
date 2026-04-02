package core_logger // пакет для логгера приложения

import (
	"context"
	"fmt"           // форматирование строк и ошибок
	"os"            // работа с файловой системой
	"path/filepath" // безопасная работа с путями
	"time"          // работа со временем

	"go.uber.org/zap"         // основной логгер
	"go.uber.org/zap/zapcore" // низкоуровневые настройки zap
)

// Logger — обёртка над zap.Logger
type Logger struct {
	*zap.Logger // embedding: все методы zap.Logger доступны напрямую

	file *os.File // файл, в который пишутся логи (нужен, чтобы потом закрыть)
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

// NewLogger — конструктор логгера
func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel() // создаём уровень логирования (можно менять на лету)

	// парсим уровень логирования из строки (например "info", "debug")
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err) // оборачиваем ошибку
	}

	// создаём папку для логов, если её нет
	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	// формируем timestamp для имени файла
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")

	// собираем полный путь до файла логов
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp), // имя файла = timestamp.log
	)

	// открываем файл для записи (создаём если нет)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	// создаём конфиг кодировщика (формат логов)
	zapConfig := zap.NewDevelopmentEncoderConfig()

	// задаём формат времени в логах
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	// создаём encoder (как будут выглядеть логи)
	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	// создаём core — куда писать логи
	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl), // вывод в консоль
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),   // вывод в файл
	)

	// создаём сам логгер + добавляем caller (файл и строка вызова)
	zapLogger := zap.New(core, zap.AddCaller())

	// возвращаем нашу обёртку
	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

// Close — закрывает файл логов
func (l *Logger) Close() {
	// закрываем файл
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err) // fallback лог
	}
}

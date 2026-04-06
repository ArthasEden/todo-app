package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger — обёртка над zap.Logger
// Встраиваем zap.Logger, чтобы использовать его методы напрямую (Info, Error и т.д.)
// + храним файл, чтобы корректно закрыть его при завершении приложения
type Logger struct {
	*zap.Logger
	file *os.File // файл, в который пишутся логи
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

// NewLogger — создаёт и настраивает логгер
func NewLogger(config Config) (*Logger, error) {

	///////////////////////////////////////////////////////////////
	// ЧАСТЬ 1: Работа с конфигом и файловой системой
	// Здесь мы:
	// - настраиваем уровень логирования
	// - создаём папку
	// - создаём файл логов
	///////////////////////////////////////////////////////////////

	// 1. Создаём уровень логирования (обёртка над уровнем, которую можно менять в runtime)
	zapLvl := zap.NewAtomicLevel()

	// Парсим уровень из строки (например: "debug", "info", "error")
	// Если строка невалидная → будет ошибка
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	// 2. Создаём папку для логов (если её нет)
	// MkdirAll создаёт всю цепочку директорий
	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return &Logger{}, fmt.Errorf("mkdir log folder: %w", err)
	}

	// 3. Генерируем имя файла с текущим временем (чтобы каждый запуск создавал новый файл)
	timeStamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")

	// Формируем полный путь до файла логов (folder + имя файла)
	logFilePath := filepath.Join(config.Folder, fmt.Sprintf("%s.log", timeStamp))

	// Открываем (или создаём) файл для записи логов
	// O_CREATE — создать файл, если его нет
	// O_WRONLY — только запись
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	///////////////////////////////////////////////////////////////
	// ЧАСТЬ 2: Настройка и сборка логгера (zap)
	// Здесь мы:
	// - задаём формат логов
	// - указываем, куда писать (stdout + файл)
	// - создаём сам логгер
	///////////////////////////////////////////////////////////////

	// 4. Настраиваем формат логов (encoder config)
	// Development — человекочитаемый формат (не JSON)
	zapConfig := zap.NewDevelopmentEncoderConfig()

	// Переопределяем формат времени (по умолчанию он другой)
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	// Encoder — отвечает за внешний вид логов (формат строки)
	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	// 5. Создаём core (куда писать логи)
	// NewCore = (encoder, destination, level)
	// NewTee = писать сразу в несколько мест (мульти-вывод)
	core := zapcore.NewTee(
		// Пишем в консоль (stdout)
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),

		// Пишем в файл
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	// 6. Создаём сам логгер
	// AddCaller — добавляет файл и строку, где был вызов логгера
	zapLogger := zap.New(core, zap.AddCaller())

	// Возвращаем нашу обёртку (с zap.Logger и файлом)
	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		file:   l.file,
	}
}

// Close — закрывает файл логов
// Важно вызывать при завершении приложения
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}

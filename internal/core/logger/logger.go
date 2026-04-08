// Package core_logger предоставляет структурированный логгер на базе go.uber.org/zap.
// Логи пишутся одновременно в stdout и в файл.
//
// Логгер передаётся через context.Context (паттерн «logger in context»),
// что позволяет автоматически добавлять к каждому сообщению
// request_id и другие поля, установленные в middleware.
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

// loggerContextKey — приватный тип ключа для context.WithValue.
// Использование отдельного типа (а не string) исключает коллизии ключей
// с другими пакетами, которые тоже хранят данные в контексте.
type loggerContextKey struct{}

var (
	key = loggerContextKey{}
)

// Logger — обёртка над *zap.Logger, которая дополнительно хранит
// файловый дескриптор для корректного закрытия при завершении приложения.
type Logger struct {
	*zap.Logger

	file *os.File
}

// ToContext кладёт логгер в контекст. Вызывается в middleware Logger,
// чтобы все последующие обработчики могли получить логгер с request_id.
func ToContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, key, log)
}

// FromContext извлекает логгер из контекста.
// Паникует, если логгер не был добавлен — это программная ошибка,
// означающая, что middleware Logger не был подключён.
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(key).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

// NewLogger создаёт логгер, который пишет в stdout и в файл одновременно.
// Каждый запуск создаёт новый лог-файл с именем вида "2006-01-02T15-04-05.000000.log".
//
// zapcore.NewTee объединяет несколько «ядер» (outputs) в одно:
// запись в одно ядро — автоматически запись во все.
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

// With создаёт дочерний логгер с дополнительными полями.
// Переопределяем метод, чтобы возвращать *core_logger.Logger (с файлом),
// а не базовый *zap.Logger.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		file:   l.file,
	}
}

// Close закрывает файл логов. Должен вызываться через defer в main().
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}

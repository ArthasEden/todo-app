// Package domain содержит доменные модели и бизнес-логику приложения.
// Это самый внутренний слой архитектуры — он не зависит ни от какого другого пакета проекта.
package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/ArthasEden/todo-app/internal/core/errors"
)

// User — доменная сущность пользователя.
//
// PhoneNumber — nil означает отсутствие номера (NULL в базе данных).
// Version — счётчик для оптимистичной блокировки: см. Task.Version.
type User struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

// NewUser — конструктор для восстановления пользователя по имеющему набору данных
func NewUser(id, version int, fullName string, phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(UninitializedID, UninitializedVersion, fullName, phoneNumber)
}

// Validate проверяет инварианты пользователя.
// Формат телефона: начинается с «+», далее только цифры, длина 10–15 символов.
// Пример: +79001234567
func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf("invalid `FullName` len: %d: %w", fullNameLength, core_errors.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("invalid `PhoneNumber` len: %d: %w", phoneNumberLen, core_errors.ErrInvalidArgument)
		}
	}

	// regexp.MustCompile паникует при невалидном паттерне — это допустимо,
	// так как паттерн — константа, известная на этапе компиляции.
	// В продакшн-коде регулярное выражение лучше вынести в переменную пакета.
	re := regexp.MustCompile(`^\+[0-9]+$`)

	if !re.MatchString(*u.PhoneNumber) {
		return fmt.Errorf("invalid `PhoneNumber` format: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

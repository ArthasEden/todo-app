package domain

import (
	"fmt"
	"regexp"
	"todoapp/internal/core/sentinels"
	"uuid"
)

var defVersion = 1

type User struct {
	ID          uuid.UUID
	Version     int
	FullName    string
	PhoneNumber *string
}

// NewUser — конструктор для восстановления пользователя по имеющему набору данных
func NewUser(
	id uuid.UUID,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

// CreateUser создаёт нового пользователя с автоматически сгенерированными
// ID (UUID v4) и начальной версией 1.
func CreateUser(
	fullName string,
	phoneNumber *string,
) User {
	var (
		id      = uuid.New()
		version = defVersion
	)

	return NewUser(
		id,
		version,
		fullName,
		phoneNumber,
	)
}

// Validate проверяет инварианты пользователя.
// Формат телефона: начинается с «+», далее только цифры, длина 10–15 символов.
// Пример: +79001234567
func (u *User) Validate() error {
	fullNameLen := len([]rune(u.FullName))
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLen,
			sentinels.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLen,
				sentinels.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				sentinels.ErrInvalidArgument,
			)
		}
	}

	return nil
}

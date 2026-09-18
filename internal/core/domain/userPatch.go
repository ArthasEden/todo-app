package domain

import (
	"fmt"
	"todoapp/internal/core/sentinels"
)

// UserPatch содержит изменения для частичного обновления пользователя (PATCH).
// Каждое поле обёрнуто в Nullable, чтобы различать «не передано» и «передано null».
// Подробнее о Nullable: см. internal/core/domain/nullable.go.
type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

// NewUserPatch — конструктор UserPatch.
func NewUserPatch(
	fullName Nullable[string],
	phoneNumber Nullable[string],
) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

// Validate проверяет корректность патча до его применения.
// FullName является обязательным полем и не может быть обнулён.
func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"`FullName` can't be patched to NULL: %w",
			sentinels.ErrInvalidArgument,
		)
	}

	return nil
}

// ApplyPatch применяет изменения к пользователю.
// Используется та же техника «копия → изменение → валидация → замена», что и в Task.ApplyPatch.
func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}

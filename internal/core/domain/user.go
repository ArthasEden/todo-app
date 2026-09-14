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

func NewUser(
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          uuid.New(),
		Version:     defVersion,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

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

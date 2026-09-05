package domain

import "uuid"

type User struct {
	id      uuid.UUID
	version int

	fullName    string
	phoneNumber *string
}

package users_transport_http

import (
	"encoding/json"
	"net/http"
	"uuid"
)

type CreateUserReq struct {
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

type CreateUserRes struct {
	ID      uuid.UUID `json:"id"`
	Vesrion int       `json:"version"`

	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func (u *UsersHTTPHandler) CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var dtoReq CreateUserReq

	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {

		return
	}
}

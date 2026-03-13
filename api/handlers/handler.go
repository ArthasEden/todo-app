package handlers

import (
	"todo-app/service"
)

type Handler struct {
	list *service.List
}

func NewHandler(list *service.List) *Handler {
	return &Handler{
		list: list,
	}
}

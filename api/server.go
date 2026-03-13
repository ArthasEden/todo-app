package api

import (
	"todo-app/api/handlers"
)

type Server struct {
	handlers *handlers.Handler
}

func NewServer(handlers *handlers.Handler) *Server {
	return &Server{
		handlers: handlers,
	}
}

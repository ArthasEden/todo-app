package api

import (
	"net/http"
	"todo-app/api/handlers"

	"github.com/gorilla/mux"
)

var (
	tasks      = "/tasks"
	tasksTitle = "/tasks/{title}"
)

type Server struct {
	handlers *handlers.Handler
}

func NewServer(handlers *handlers.Handler) *Server {
	return &Server{
		handlers: handlers,
	}
}

func (s *Server) Run() error {
	r := mux.NewRouter()

	r.Path(tasks).Methods("POST").HandlerFunc(s.handlers.Post)
	r.Path(tasks).Methods("GET").HandlerFunc(s.handlers.GetAll)
	r.Path(tasksTitle).Methods("PATCH").HandlerFunc(s.handlers.Patch)
	r.Path(tasksTitle).Methods("DELETE").HandlerFunc(s.handlers.Delete)

	return http.ListenAndServe(":9091", r)
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"todo-app/api/dto"
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

func SendJSON(w http.ResponseWriter, data any, statusCode int) {
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		http.Error(w, dto.ErrDTO(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(b); err != nil {
		log.Println("couldn't write response:", err)
		return
	}
}

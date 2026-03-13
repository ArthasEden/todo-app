package handlers

import (
	"errors"
	"net/http"
	"todo-app/api/dto"
	"todo-app/service"

	"github.com/gorilla/mux"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	if err := h.list.Delete(title); err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			http.Error(w, dto.ErrDTO(err), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

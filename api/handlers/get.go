package handlers

import (
	"errors"
	"net/http"
	"todo-app/api/dto"
	"todo-app/service"

	"github.com/gorilla/mux"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var (
		data any
		err  error
	)

	title := mux.Vars(r)["title"]

	switch title {
	case "all":
		data = h.list.GetAll()
	default:
		data, err = h.list.GetOne(title)
	}

	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			http.Error(w, dto.ErrDTO(err), http.StatusNotFound)
			return
		}
		http.Error(w, dto.ErrDTO(err), http.StatusBadRequest)
		return
	}

	SendJSON(w, data, http.StatusOK)
}

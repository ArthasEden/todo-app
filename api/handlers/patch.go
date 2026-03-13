package handlers

import (
	"errors"
	"net/http"
	"todo-app/api/dto"
	"todo-app/service"

	"github.com/gorilla/mux"
)

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	var (
		title  = mux.Vars(r)["title"]
		status = r.URL.Query().Get("completed")
		err    error
		task   service.Task
	)

	switch status {
	case "true":
		task, err = h.list.Complete(title)
	case "false":
		task, err = h.list.Uncomplete(title)
	default:
		http.Error(w, dto.ErrDTO(errors.New("invalid completed parameter")), http.StatusBadRequest)
		return
	}

	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			http.Error(w, dto.ErrDTO(err), http.StatusNotFound)
			return
		}
		http.Error(w, dto.ErrDTO(err), http.StatusInternalServerError)
		return
	}

	SendJSON(w, task, http.StatusOK)
}

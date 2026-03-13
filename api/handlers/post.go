package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"todo-app/api/dto"
	"todo-app/service"
)

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	tDTO := &dto.TDTO{}

	if err := json.NewDecoder(r.Body).Decode(tDTO); err != nil {
		http.Error(w, dto.ErrDTO(err), http.StatusBadRequest)
		return
	}

	task := service.NewTask(tDTO.Title, tDTO.Text)

	if err := h.list.Add(*task); err != nil {
		if errors.Is(err, service.ErrTaskAlreadyExist) {
			http.Error(w, dto.ErrDTO(err), http.StatusConflict)
			return
		}
		http.Error(w, dto.ErrDTO(err), http.StatusInternalServerError)
		return
	}

	SendJSON(w, task, http.StatusCreated)
}

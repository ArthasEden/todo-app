package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"todoapp/internal/core/sentinels"
)

type ErrorResponse struct {
	Error   string `json:"error"   example:"full error text"`
	Message string `json:"message" example:"short human-readable message"`
}

func (w *ResponseWriter) JSONResponse(response any, statusCode int) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.log.Error("write HTTP response", slog.Any("error", err))
	}
}

func (w *ResponseWriter) NoContentResponse() {
	w.WriteHeader(http.StatusNoContent)
}

func (w *ResponseWriter) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...any)
	)

	switch {
	case errors.Is(err, sentinels.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = w.log.Info

	case errors.Is(err, sentinels.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = w.log.Debug

	case errors.Is(err, sentinels.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = w.log.Info

	default:
		statusCode = http.StatusInternalServerError
		logFunc = w.log.Error
	}

	logFunc(msg, err.Error())

	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}
	w.JSONResponse(response, statusCode)
}

func (w *ResponseWriter) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	w.log.Error(msg, slog.Any("error", err))

	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}
	w.JSONResponse(response, statusCode)
}

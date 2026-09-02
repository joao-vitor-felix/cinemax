package controller

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)



type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, d any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(d)
}

func DecodeJSON(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

func WriteError(w http.ResponseWriter, status int, code, msg string) {
	WriteJSON(w, status, ErrorResponse{
		Code:    code,
		Message: msg,
	})
}

func HandleError(w http.ResponseWriter, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		slog.Error("application error", "code", appErr.Code, "message", appErr.Message)
		WriteError(w, appErr.StatusCode, appErr.Code, appErr.Message)
		return
	}
	WriteError(
		w,
		http.StatusInternalServerError,
		domain.InternalServerError.Code,
		domain.InternalServerError.Message,
	)
}

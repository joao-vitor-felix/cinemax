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
	if d == nil {
		w.WriteHeader(status)
		return
	}
	if err := json.NewEncoder(w).Encode(d); err != nil {
		slog.Error("failed to encode json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Code:    domain.InternalServerError.Code,
			Message: domain.InternalServerError.Message,
		})
		return
	}
	w.WriteHeader(status)
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
	slog.Error("unexpected error", "error", err)
	WriteError(
		w,
		http.StatusInternalServerError,
		domain.InternalServerError.Code,
		domain.InternalServerError.Message,
	)
}

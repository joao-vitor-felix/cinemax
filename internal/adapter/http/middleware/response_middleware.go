package middleware

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteErrorResponse(w http.ResponseWriter, err *domain.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	resp := ErrorResponse{
		Code:    err.Code,
		Message: err.Message,
	}
	jsonBytes, _ := json.Marshal(resp)
	w.Write(jsonBytes)
}

type AppHandler func(w http.ResponseWriter, r *http.Request) (map[string]any, error)

func MakeHandler(fn AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := fn(w, r)

		if err != nil {
			var appErr *domain.AppError
			if errors.As(err, &appErr) {
				//TODO: make logging better
				slog.Error("App Error: " + appErr.Code + " - " + appErr.Message)
				WriteErrorResponse(w, appErr)
				return
			}
			slog.Error("Unhandled error", slog.Any("err", err))
			WriteErrorResponse(w, domain.InternalServerError)
			return
		}

		data := res["res"]
		status := res["status"].(int)

		if data == nil {
			w.WriteHeader(status)
			return
		}

		jsonBytes, err := json.Marshal(data)
		if err != nil {
			slog.Error("Failed to marshal JSON response", slog.Any("err", err))
			WriteErrorResponse(w, domain.InternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(jsonBytes)
	}
}

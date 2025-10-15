package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var defaultCors cors.Options = cors.Options{
	AllowedOrigins: []string{
		//TODO: remove it later and place into env variables
		"*",
		"http://localhost:3000",
	},
	AllowedHeaders: []string{
		"Content-Type",
		"Authorization",
		"Cache-Control",
	},
}

func SetupRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/"))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(cors.Handler(defaultCors))
	// use it later on routes that should not be cached
	// r.Use(middleware.NoCache)
	r.Get("/placeholder", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("testado"))
	})
	return r
}

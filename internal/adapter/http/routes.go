package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	m "github.com/joao-vitor-felix/cinemax/internal/adapter/http/middleware"
	"github.com/joao-vitor-felix/cinemax/internal/factory"
	"github.com/joao-vitor-felix/cinemax/internal/locale"
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

func SetupRoutes(container *factory.Container) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/"))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(cors.Handler(defaultCors))

	bundle := locale.NewBundle()
	r.Use(m.LocalizeMiddleware(bundle))
	// use it later on routes that should not be cached
	// r.Use(middleware.NoCache)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", m.MakeHandler(container.UserController.Register))
	})
	return r
}

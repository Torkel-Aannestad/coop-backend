package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(app.panicRecovery)
	r.Get("/v1/healthcheck", app.healthcheckHandler)

	return r
}

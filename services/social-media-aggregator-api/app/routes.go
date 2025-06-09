package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/v1/healthcheck", app.healthcheckHandler)
	r.Get("/v1/messages", app.latestMessagesHandler)
	r.Post("/v1/messages", app.createMessageHandler)

	return r
}

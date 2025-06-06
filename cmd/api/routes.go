package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("working"))
}

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/v1/healthcheck", app.healthcheckHandler)

	return r
}

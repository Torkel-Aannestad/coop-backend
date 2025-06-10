package main

import (
	"net/http"

	"github.com/Torkel-Aannestad/coop-backend/services/social-media-aggregator-api/internal/database"
)

func (app *application) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ExternalId string `json:"external_id"`
		Author     string `json:"author"`
		Body       string `json:"body"`
		Platform   string `json:"platform"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	message := database.Message{
		ExternalId: input.ExternalId,
		Author:     input.Author,
		Body:       input.Body,
		Platform:   input.Platform,
	}

	// validate input

	err = app.models.Messages.Insert(&message)
	if err != nil {
		//todo handle db unique error for enternalId
		app.serverErrorResponse(w, r, err)
		return
	}
	app.logger.Info("message created", "ExternalId", message.ExternalId) //remove

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": message}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) latestMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// returns the 25 last messages from db

	messages, err := app.models.Messages.GetList(10, 0)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, 200, envelope{"messages": messages}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

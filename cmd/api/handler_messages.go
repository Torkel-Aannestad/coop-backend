package main

import "net/http"

func (app *application) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("createMessageHandler"))
}

func (app *application) listMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("listMessageHandler"))
}

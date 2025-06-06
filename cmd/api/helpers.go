package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJson(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(js))

	return nil
}

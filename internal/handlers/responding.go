package handlers

import (
	"encoding/json"
	"net/http"
)

type Failure struct {
	Description string `json:"description"`
	Error       string `json:"error,omitempty"`
}

var EmptyJSON = map[string]interface{}{}

func respond(w http.ResponseWriter, code int, response interface{}) {
	body, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(code)
	w.Write(body)
}

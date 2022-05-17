package api

import (
	"encoding/json"
	"net/http"
)

type ResponseMessage struct {
	Message string
}

func respond(w http.ResponseWriter, code int, out []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(out); err != nil {
		panic(err)
	}
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	respond(w, code, resp)
}

func RespondWithError(w http.ResponseWriter, code int, payload error) {
	resp, err := json.Marshal(ResponseMessage{Message: payload.Error()})
	if err != nil {
		panic(err)
	}
	respond(w, code, resp)
}

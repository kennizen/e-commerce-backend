package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResUser struct {
	Msg string `json:"message"`
}

func NewResUser(msg string) *ResUser {
	return &ResUser{
		Msg: msg,
	}
}

func SendResp(msg string, code int, w http.ResponseWriter) {
	data := NewResUser(msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		log.Fatal("Error parsing JSON")
	}
}

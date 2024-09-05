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

func SendMsg(msg string, status int, w http.ResponseWriter) {
	data := NewResUser(msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		log.Fatal("Error parsing JSON")
	}
}

func SendJson(resp any, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(resp)

	if err != nil {
		log.Fatal("Error parsing JSON")
	}
}

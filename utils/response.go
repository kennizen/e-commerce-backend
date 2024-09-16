package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResUser struct {
	Msg string `json:"message"`
}

type ResUserWithData struct {
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func NewResUser(msg string) *ResUser {
	return &ResUser{
		Msg: msg,
	}
}

func NewResUserWithData(msg string, data any) *ResUserWithData {
	return &ResUserWithData{
		Msg:  msg,
		Data: data,
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

func SendJson(payload ResUserWithData, status int, w http.ResponseWriter) {
	data := NewResUserWithData(payload.Msg, payload.Data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		log.Fatal("Error parsing JSON")
	}
}

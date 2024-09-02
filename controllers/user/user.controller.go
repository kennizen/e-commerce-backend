package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	service "github.com/kennizen/e-commerce-backend/services/user"
)

func LoginController(w http.ResponseWriter, r *http.Request) {
	var user service.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("user", user)

	user.LoginUserService(w)
}

func RegisterController(w http.ResponseWriter, r *http.Request) {
	var user service.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("user", user)

	user.RegisterUserService(w)
}

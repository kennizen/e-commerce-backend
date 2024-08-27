package routes

import (
	"net/http"

	"github.com/kennizen/e-commerce-backend/controllers/user"
)

func RegisterRoutes(router *http.ServeMux) {
	// auth routes
	router.HandleFunc("POST /login", controller.LoginController)
	router.HandleFunc("POST /register", controller.RegisterController)
}

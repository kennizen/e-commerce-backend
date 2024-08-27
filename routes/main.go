package routes

import (
	"net/http"

	"github.com/kennizen/e-commerce-backend/controllers"
)

func RegisterRoutes(router *http.ServeMux) {
	// auth routes
	router.HandleFunc("POST /login", controllers.LoginController)
}

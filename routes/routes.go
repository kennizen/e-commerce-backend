package routes

import (
	"net/http"

	controller "github.com/kennizen/e-commerce-backend/controllers"
	// "github.com/kennizen/e-commerce-backend/middlewares"
)

func RegisterRoutes(router *http.ServeMux) {
	// auth routes
	router.HandleFunc("POST /login", controller.LoginController)
	router.HandleFunc("POST /register", controller.RegisterController)

	// product routes
	router.HandleFunc("GET /products", controller.GetProducts)
	router.HandleFunc("GET /product/{id}", controller.GetProduct)
}

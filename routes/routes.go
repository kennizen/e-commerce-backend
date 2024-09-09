package routes

import (
	"fmt"
	"net/http"
	"os"

	controller "github.com/kennizen/e-commerce-backend/controllers"
	_ "github.com/kennizen/e-commerce-backend/docs"
	"github.com/kennizen/e-commerce-backend/middlewares"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title E-Commerce Backend API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func RegisterRoutes(router *http.ServeMux) {
	// api docs
	router.HandleFunc("GET /swagger/*", httpSwagger.Handler(httpSwagger.URL(
		fmt.Sprintf("http://localhost:%s/swagger/doc.json", os.Getenv("API_PORT")),
	)))

	// auth routes
	router.HandleFunc("POST /login", controller.LoginController)
	router.HandleFunc("POST /register", controller.RegisterController)

	// product routes
	router.HandleFunc("GET /products", controller.GetProducts)
	router.HandleFunc("GET /product/{productId}", controller.GetProduct)
	router.HandleFunc("POST /product/favorite/{productId}", middlewares.Authenticate(controller.MarkFavorite))
	router.HandleFunc("DELETE /product/favorite/{productId}", middlewares.Authenticate(controller.UnMarkFavorite))
	router.HandleFunc("GET /product/favorites", middlewares.Authenticate(controller.GetFavorites))
	router.HandleFunc("POST /product/review", middlewares.Authenticate(controller.AddProductReview))

	// cart routes
	router.HandleFunc("GET /cart", middlewares.Authenticate(controller.GetCart))
	router.HandleFunc("POST /cart", middlewares.Authenticate(controller.AddToCart))
	router.HandleFunc("DELETE /cart/{productId}", middlewares.Authenticate(controller.RemoveFromCart))
	router.HandleFunc("PATCH /cart", middlewares.Authenticate(controller.UpdateCartItems))
}

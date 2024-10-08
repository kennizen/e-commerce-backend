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
// @description This is a dummy backend for an ecommerce store.
// @termsOfService http://swagger.io/terms/

// @contact.email prachurjyagogoi123@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
func RegisterRoutes(router *http.ServeMux) {
	// api docs
	router.HandleFunc("GET /swagger/*", httpSwagger.Handler(httpSwagger.URL(
		fmt.Sprintf("http://localhost:%s/swagger/doc.json", os.Getenv("API_PORT")),
	)))

	// auth routes
	router.HandleFunc("POST /login", controller.LoginController)
	router.HandleFunc("POST /register", controller.RegisterController)
	router.HandleFunc("GET /renew-access-token", controller.RenewAccessToken)

	// user routes
	router.HandleFunc("PUT /user", middlewares.Authenticate(controller.UpdateUserDetails))
	router.HandleFunc("DELETE /user", middlewares.Authenticate(controller.DeleteUser))
	router.HandleFunc("POST /user/address", middlewares.Authenticate(controller.AddAddress))
	router.HandleFunc("PUT /user/address/{addressId}", middlewares.Authenticate(controller.UpdateAddress))
	router.HandleFunc("DELETE /user/address/{addressId}", middlewares.Authenticate(controller.DeleteAddress))
	router.HandleFunc("GET /user/addresses", middlewares.Authenticate(controller.GetAddresses))

	// product routes
	router.HandleFunc("GET /products", controller.GetProducts)
	router.HandleFunc("GET /product/{productId}", controller.GetProduct)
	router.HandleFunc("POST /favorite/product/{productId}", middlewares.Authenticate(controller.MarkFavorite))
	router.HandleFunc("DELETE /favorite/product/{productId}", middlewares.Authenticate(controller.UnMarkFavorite))
	router.HandleFunc("GET /favorite/products", middlewares.Authenticate(controller.GetFavorites))
	router.HandleFunc("POST /review/product/{productId}", middlewares.Authenticate(controller.AddProductReview))
	router.HandleFunc("PUT /review/product/{reviewId}", middlewares.Authenticate(controller.UpdateProductReview))
	router.HandleFunc("DELETE /review/product/{reviewId}", middlewares.Authenticate(controller.DeleteProductReview))
	router.HandleFunc("GET /product/{productId}/reviews", controller.GetProductReviewsByProductId)

	// cart routes
	router.HandleFunc("GET /cart", middlewares.Authenticate(controller.GetCart))
	router.HandleFunc("POST /cart/product/{productId}", middlewares.Authenticate(controller.AddToCart))
	router.HandleFunc("DELETE /cart/product/{productId}", middlewares.Authenticate(controller.RemoveFromCart))
	router.HandleFunc("PUT /cart/product/{productId}", middlewares.Authenticate(controller.UpdateCartItems))

	// order routes
	router.HandleFunc("POST /order", middlewares.Authenticate(controller.PlaceOrder))
	router.HandleFunc("GET /orders", middlewares.Authenticate(controller.GetOrders))
}

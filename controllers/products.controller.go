package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kennizen/e-commerce-backend/middlewares"
	service "github.com/kennizen/e-commerce-backend/services"
	"github.com/kennizen/e-commerce-backend/utils"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting all products")

	queryParams := r.URL.Query()

	var page, limit string

	if queryParams.Get("page") == "" {
		page = "1"
	} else {
		page = queryParams.Get("page")
	}

	if queryParams.Get("limit") == "" {
		limit = "-1"
	} else {
		limit = queryParams.Get("limit")
	}

	page1, err := strconv.Atoi(page)
	limit1, err1 := strconv.Atoi(limit)

	if err != nil || err1 != nil {
		log.Fatalln("Error converting query params")
	}

	service.GetProducts(page1, limit1, w)
}

// ---------------------------------------------------------------------------------------- //

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("productId")

	productId, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("Invalid product id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	service.GetProduct(productId, w)
}

// ---------------------------------------------------------------------------------------- //

func MarkFavorite(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	id := r.PathValue("productId")

	_, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("Invalid product id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	service.MarkFavorite(userId.(string), id, w)
}

// ---------------------------------------------------------------------------------------- //

func UnMarkFavorite(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	id := r.PathValue("productId")

	_, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("Invalid product id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	service.UnMarkFavorite(userId.(string), id, w)
}

// ---------------------------------------------------------------------------------------- //

func GetFavorites(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	service.GetFavorites(userId.(string), w)
}

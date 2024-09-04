package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		limit = "10"
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

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	productId, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("Invalid product id")
		utils.SendError("Invalid id", http.StatusBadRequest, w)
		return
	}

	service.GetProduct(productId, w)
}

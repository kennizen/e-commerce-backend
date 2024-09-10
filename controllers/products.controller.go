package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kennizen/e-commerce-backend/middlewares"
	service "github.com/kennizen/e-commerce-backend/services"
	"github.com/kennizen/e-commerce-backend/utils"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
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

// ---------------------------------------------------------------------------------------- //

func AddProductReview(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	productId := r.PathValue("productId")

	_, strConvErr := strconv.Atoi(productId)

	if strConvErr != nil {
		fmt.Println("Invalid product id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	var payload struct {
		Review string  `validate:"required"`
		Rating float32 `validate:"required"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		fmt.Println("Error in json decoding", err.Error())
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	valErr := utils.Validate(payload)

	if valErr != nil {
		fmt.Println("Invalid payload", valErr.Error())
		utils.SendMsg("Invalid payload", http.StatusBadRequest, w)
		return
	}

	service.AddProductReview(service.ProductReviewArgs{
		Review:    payload.Review,
		Rating:    payload.Rating,
		ProductId: productId,
		UserId:    userId.(string),
	}, w)
}

// ---------------------------------------------------------------------------------------- //

func UpdateProductReview(w http.ResponseWriter, r *http.Request) {
	reviewId := r.PathValue("reviewId")

	_, strConvErr := strconv.Atoi(reviewId)

	if strConvErr != nil {
		fmt.Println("Invalid review id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	var payload struct {
		Review string  `validate:"required"`
		Rating float32 `validate:"required"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		fmt.Println("Error in json decoding", err.Error())
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	valErr := utils.Validate(payload)

	if valErr != nil {
		fmt.Println("Invalid payload", valErr.Error())
		utils.SendMsg("Invalid payload", http.StatusBadRequest, w)
		return
	}

	service.UpdateProductReview(service.ProductUpdateArgs{
		ReviewId: reviewId,
		Review:   payload.Review,
		Rating:   payload.Rating,
	}, w)
}

// ---------------------------------------------------------------------------------------- //

func DeleteProductReview(w http.ResponseWriter, r *http.Request) {
	reviewId := r.PathValue("reviewId")

	_, strConvErr := strconv.Atoi(reviewId)

	if strConvErr != nil {
		fmt.Println("Invalid review id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	service.DeleteProductReview(reviewId, w)
}

// ---------------------------------------------------------------------------------------- //

func GetProductReviewsByProductId(w http.ResponseWriter, r *http.Request) {
	productId := r.PathValue("productId")

	_, strConvErr := strconv.Atoi(productId)

	if strConvErr != nil {
		fmt.Println("Invalid product id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	service.GetProductReviewsByProductId(productId, w)
}

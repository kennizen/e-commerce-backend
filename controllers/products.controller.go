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

// @Summary      Get all products.
// @Description  API for fetching all products.
// @Tags         Products
// @Produce      json
// @Param        page query int true "Page number"
// @Param        limit query int true "The number of products to fetch"
// @Success      200  {object} utils.ResUserWithData{data=service.ProductsResponse}
// @Failure      500  {object} utils.ResUser
// @Router       /products [get]
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

// @Summary      Get a product details.
// @Description  API for fetching a product.
// @Tags         Products
// @Produce      json
// @Param        productId path string true "Product ID"
// @Success      200  {object} utils.ResUserWithData{data=models.Product}
// @Failure      400  {object} utils.ResUser
// @Failure      404  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /product/{productId} [get]
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

// @Summary      Mark a product favorite.
// @Description  API for making a product as favorite.
// @Tags         Products
// @Produce      json
// @Param        productId path string true "Product ID"
// @Success      200  {object} utils.ResUser
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      409  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /favorite/product/{productId} [post]
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

// @Summary      UnMark a product as favorite.
// @Description  API for removing a product as favorite.
// @Tags         Products
// @Produce      json
// @Param        productId path string true "Product ID"
// @Success      200  {object} utils.ResUser
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /favorite/product/{productId} [delete]
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

// @Summary      Get all favorite products.
// @Description  API for getting all favorite products.
// @Tags         Products
// @Produce      json
// @Success      200  {object} utils.ResUserWithData{data=[]models.Product}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /favorite/products [get]
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

// @Summary      Add a product review.
// @Description  API for adding a product review by user.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        user body service.ProductReviewPayload true "Product review payload"
// @Param        productId path string true "Product id"
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=models.ProductReview}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /review/product/{productId} [post]
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

	var payload service.ProductReviewPayload

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

	service.AddProductReview(service.ProductReviewPayload{
		Review: payload.Review,
		Rating: payload.Rating,
	}, productId, userId.(string), w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Update a product review.
// @Description  API for updating a product review by user.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        user body service.ProductReviewPayload true "Product review payload"
// @Param        reviewId path string true "Review ID"
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=models.ProductReview}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /review/product/{reviewId} [put]
func UpdateProductReview(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	reviewId := r.PathValue("reviewId")

	_, strConvErr := strconv.Atoi(reviewId)

	if strConvErr != nil {
		fmt.Println("Invalid review id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	var payload service.ProductReviewPayload

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

	service.UpdateProductReview(service.ProductReviewPayload{
		Review: payload.Review,
		Rating: payload.Rating,
	}, reviewId, userId.(string), w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Delete a product review.
// @Description  API for deleting a product review by user.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        reviewId path string true "Review ID"
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=models.ProductReview}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /review/product/{reviewId} [delete]
func DeleteProductReview(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	reviewId := r.PathValue("reviewId")

	_, strConvErr := strconv.Atoi(reviewId)

	if strConvErr != nil {
		fmt.Println("Invalid review id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	service.DeleteProductReview(reviewId, userId.(string), w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Get all reviews for a product.
// @Description  API for fetching all the reviews for a product made by multiple users.
// @Tags         Products
// @Produce      json
// @Param        productId path string true "Product ID"
// @Success      200  {object} utils.ResUserWithData{data=[]service.AllReviewsResponse}
// @Failure      400  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /product/{productId}/reviews [get]
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

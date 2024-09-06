package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kennizen/e-commerce-backend/middlewares"
	service "github.com/kennizen/e-commerce-backend/services"
	"github.com/kennizen/e-commerce-backend/utils"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	var product struct {
		Id       string
		Quantity int
	}

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		fmt.Println("Error in json decoding", err.Error())
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	service.AddToCart(service.CartArgs{
		UserId:    userId.(string),
		ProductId: product.Id,
		Quantity:  product.Quantity,
	},
		w,
	)
}

// ---------------------------------------------------------------------------------------- //

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
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

	service.RemoveFromCart(userId.(string), id, w)
}

// ---------------------------------------------------------------------------------------- //

func UpdateCartItems(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	var product struct {
		Id       string
		Quantity int
	}

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		fmt.Println("Error in json decoding", err.Error())
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	service.UpdateCartItems(service.CartArgs{
		UserId:    userId.(string),
		ProductId: product.Id,
		Quantity:  product.Quantity,
	},
		w,
	)
}

// ---------------------------------------------------------------------------------------- //

func GetCart(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	service.GetCart(userId.(string), w)
}

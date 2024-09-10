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

	productId := r.PathValue("productId")

	_, err := strconv.Atoi(productId)

	if err != nil {
		fmt.Println("Invalid product id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	var payload struct {
		Quantity int `validate:"required,gte=1"`
	}

	err1 := json.NewDecoder(r.Body).Decode(&payload)

	if err1 != nil {
		fmt.Println("Error in json decoding", err1.Error())
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	valErr := utils.Validate(payload)

	if valErr != nil {
		fmt.Println("Invalid payload", valErr.Error())
		utils.SendMsg("Invalid payload", http.StatusBadRequest, w)
		return
	}

	service.AddToCart(service.CartArgs{
		UserId:    userId.(string),
		ProductId: productId,
		Quantity:  payload.Quantity,
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

	id := r.PathValue("productId")

	_, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("Invalid product id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	var payload struct {
		Quantity int `validate:"required,gte=1"`
	}

	err1 := json.NewDecoder(r.Body).Decode(&payload)

	if err1 != nil {
		fmt.Println("Error in json decoding", err1.Error())
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	service.UpdateCartItems(service.CartArgs{
		UserId:    userId.(string),
		ProductId: id,
		Quantity:  payload.Quantity,
	}, w)
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

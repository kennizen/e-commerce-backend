package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	service.AddToCart(service.AddToCartArgs{
		UserId:    userId.(string),
		ProductId: product.Id,
		Quantity:  product.Quantity,
	},
		w,
	)
}

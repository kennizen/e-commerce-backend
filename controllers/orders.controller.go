package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kennizen/e-commerce-backend/middlewares"
	service "github.com/kennizen/e-commerce-backend/services"
	"github.com/kennizen/e-commerce-backend/utils"
)

// @Summary      Place an order.
// @Description  API for placing an order with the products in the cart.
// @Tags         Order
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer accessToken"
// @Param        order body service.OrdersPayload true "Order payload"
// @Success      200  {object} utils.ResUser
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /order [post]
func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	var payload service.OrdersPayload

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

	service.PlaceOrder(service.OrdersPayload{
		Products:      payload.Products,
		AddressUsed:   payload.AddressUsed,
		PaymentMethod: payload.PaymentMethod,
	}, userId.(string), w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Get all orders of a user.
// @Description  API for fetching all the orders of a user.
// @Tags         Order
// @Produce      json
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=[]service.AllOrdersResponse}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /orders [get]
func GetOrders(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	service.GetOrders(userId.(string), w)
}

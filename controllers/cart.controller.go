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

// @Summary      Add product to cart.
// @Description  API for adding an product the cart.
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer accessToken"
// @Param        productId path string true "Product ID"
// @Param        order body service.AddToCartPayload true "Cart payload"
// @Success      200  {object} utils.ResUser
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /cart/product/{productId} [post]
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

	var payload service.AddToCartPayload

	err1 := json.NewDecoder(r.Body).Decode(&payload)

	if err1 != nil {
		fmt.Println("Error in json decoding", err1.Error())
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	valErr := utils.Validate(payload)

	if valErr != nil {
		fmt.Println("Invalid payload", valErr.Error())
		utils.SendMsg(valErr.Error(), http.StatusBadRequest, w)
		return
	}

	str, err2 := service.AddToCart(service.AddToCartPayload{
		Quantity: payload.Quantity,
	},
		userId.(string),
		productId,
	)

	if err2 != nil {
		utils.SendMsg(err2.(*utils.HttpError).Message, err2.(*utils.HttpError).Status, w)
		return
	}

	utils.SendMsg(str, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Remove product from cart.
// @Description  API for removing a product from cart.
// @Tags         Cart
// @Produce      json
// @Param        Authorization header string true "Bearer accessToken"
// @Param        productId path string true "Product ID"
// @Success      200  {object} utils.ResUser
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /cart/product/{productId} [delete]
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

	str, err1 := service.RemoveFromCart(userId.(string), id)

	if err1 != nil {
		utils.SendMsg(err1.(*utils.HttpError).Message, err1.(*utils.HttpError).Status, w)
		return
	}

	utils.SendMsg(str, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Update a cart product.
// @Description  API for updating the quantity for a cart product.
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer accessToken"
// @Param        productId path string true "Product ID"
// @Param        order body service.AddToCartPayload true "Cart payload"
// @Success      200  {object} utils.ResUser
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /cart/product/{productId} [put]
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

	var payload service.AddToCartPayload

	err1 := json.NewDecoder(r.Body).Decode(&payload)

	if err1 != nil {
		fmt.Println("Error in json decoding", err1.Error())
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	str, err2 := service.UpdateCartItems(service.AddToCartPayload{
		Quantity: payload.Quantity,
	},
		userId.(string),
		id,
	)

	if err2 != nil {
		utils.SendMsg(err2.(*utils.HttpError).Message, err2.(*utils.HttpError).Status, w)
		return
	}

	utils.SendMsg(str, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Get all cart products.
// @Description  API for fetching all the cart products.
// @Tags         Cart
// @Produce      json
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=[]service.GetCartResponse}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /cart [get]
func GetCart(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	res, err := service.GetCart(userId.(string))

	if err != nil {
		utils.SendMsg(err.(*utils.HttpError).Message, err.(*utils.HttpError).Status, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{
		Msg:  "Cart data",
		Data: res,
	}, http.StatusOK, w)
}

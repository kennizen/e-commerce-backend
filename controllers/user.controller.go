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

func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	var payload service.UserDetailsPayload

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

	service.UpdateUserDetails(service.UserDetailsPayload{
		Firstname:  payload.Firstname,
		Middlename: payload.Middlename,
		Lastname:   payload.Lastname,
		Age:        payload.Age,
		Email:      payload.Email,
		Avatar:     payload.Avatar,
	}, userId.(string), w)
}

// ---------------------------------------------------------------------------------------- //

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	service.DeleteUser(userId.(string), w)
}

// ---------------------------------------------------------------------------------------- //

func AddAddress(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	var payload service.UserAddressPayload

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

	service.AddAddress(service.UserAddressPayload{
		Country: payload.Country,
		State:   payload.State,
		Zipcode: payload.Zipcode,
		PhoneNo: payload.PhoneNo,
		Address: payload.Address,
	}, userId.(string), w)
}

// ---------------------------------------------------------------------------------------- //

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("addressId")

	_, strErr := strconv.Atoi(id)

	if strErr != nil {
		fmt.Println("Invalid address id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	var payload service.UserAddressPayload

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

	service.UpdateAddress(service.UserAddressPayload{
		Country: payload.Country,
		State:   payload.State,
		Zipcode: payload.Zipcode,
		PhoneNo: payload.PhoneNo,
		Address: payload.Address,
	}, id, w)
}

// ---------------------------------------------------------------------------------------- //

func DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	id := r.PathValue("addressId")

	_, strErr := strconv.Atoi(id)

	if strErr != nil {
		fmt.Println("Invalid address id")
		utils.SendMsg("Invalid id", http.StatusBadRequest, w)
		return
	}

	service.DeleteAddress(id, userId.(string), w)
}

// ---------------------------------------------------------------------------------------- //

func GetAddresses(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	service.GetAddresses(userId.(string), w)
}

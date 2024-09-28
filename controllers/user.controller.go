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

// @Summary      Update user details.
// @Description  API for updating users details.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body service.UserDetailsPayload true "Update user details payload"
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=models.User}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /user [patch]
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
		utils.SendMsg(valErr.Error(), http.StatusBadRequest, w)
		return
	}

	res, err1 := service.UpdateUserDetails(service.UserDetailsPayload{
		Firstname:  payload.Firstname,
		Middlename: payload.Middlename,
		Lastname:   payload.Lastname,
		Age:        payload.Age,
		Email:      payload.Email,
		Avatar:     payload.Avatar,
	}, userId.(string))

	if err1 != nil {
		utils.SendMsg(err1.(*utils.HttpError).Message, err1.(*utils.HttpError).Status, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{
		Msg:  "User updated",
		Data: res,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Delete user.
// @Description  API for deleting a user.
// @Tags         User
// @Produce      json
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=models.User}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /user [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	res, err := service.DeleteUser(userId.(string))

	if err != nil {
		utils.SendMsg(err.(*utils.HttpError).Message, err.(*utils.HttpError).Status, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{
		Msg:  "User deleted",
		Data: res,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Add user address.
// @Description  API for adding user addresse used as delivery address.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body service.UserAddressPayload true "Add address payload"
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=models.Address}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /user/address [post]
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
		utils.SendMsg(valErr.Error(), http.StatusBadRequest, w)
		return
	}

	res, err1 := service.AddAddress(service.UserAddressPayload{
		Country: payload.Country,
		State:   payload.State,
		Zipcode: payload.Zipcode,
		PhoneNo: payload.PhoneNo,
		Address: payload.Address,
	}, userId.(string))

	if err1 != nil {
		utils.SendMsg(err1.(*utils.HttpError).Message, err1.(*utils.HttpError).Status, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{
		Msg:  "Address added",
		Data: res,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Update user address.
// @Description  API for updating a user address.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body service.UserAddressPayload true "Update address payload"
// @Param        addressId path string true "Address id"
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=models.Address}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /user/address/{addressId} [put]
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
		utils.SendMsg(valErr.Error(), http.StatusBadRequest, w)
		return
	}

	res, err1 := service.UpdateAddress(service.UserAddressPayload{
		Country: payload.Country,
		State:   payload.State,
		Zipcode: payload.Zipcode,
		PhoneNo: payload.PhoneNo,
		Address: payload.Address,
	}, id)

	if err1 != nil {
		utils.SendMsg(err1.(*utils.HttpError).Message, err1.(*utils.HttpError).Status, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{
		Msg:  "Address updated",
		Data: res,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Delete user address.
// @Description  API for deleting a user address.
// @Tags         User
// @Produce      json
// @Param        addressId path string true "Address id"
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=models.Address}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /user/address/{addressId} [delete]
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

	res, err := service.DeleteAddress(id, userId.(string))

	if err != nil {
		utils.SendMsg(err.(*utils.HttpError).Message, err.(*utils.HttpError).Status, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{
		Msg:  "Address deleted",
		Data: res,
	}, http.StatusOK, w)

}

// ---------------------------------------------------------------------------------------- //

// @Summary      Get all user addresses.
// @Description  API for fetching all user addresses.
// @Tags         User
// @Produce      json
// @Param        Authorization header string true "Bearer accessToken"
// @Success      200  {object} utils.ResUserWithData{data=[]models.Address}
// @Failure      400  {object} utils.ResUser
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /user/addresses [get]
func GetAddresses(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextKey("userID"))

	if userId == nil {
		fmt.Println("userID not found.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	res, err := service.GetAddresses(userId.(string))

	if err != nil {
		utils.SendMsg(err.(*utils.HttpError).Message, err.(*utils.HttpError).Status, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{
		Msg:  "Addresses",
		Data: res,
	}, http.StatusOK, w)
}

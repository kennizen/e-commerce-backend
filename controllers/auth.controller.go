package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	service "github.com/kennizen/e-commerce-backend/services"
	"github.com/kennizen/e-commerce-backend/utils"
)

// @Summary      User Login
// @Description  API for user login using email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body service.LoginUserPayload true "Login credentials"
// @Success      200  {object} utils.ResUserWithData{data=lib.Tokens}
// @Failure      400  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /login [post]
func LoginController(w http.ResponseWriter, r *http.Request) {
	var payload service.LoginUserPayload

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		utils.SendMsg(err.Error(), http.StatusInternalServerError, w)
		return
	}

	valErr := utils.Validate(payload)

	if valErr != nil {
		fmt.Println("Invalid payload", valErr.Error())
		utils.SendMsg("Invalid payload", http.StatusBadRequest, w)
		return
	}

	tokens, err1 := service.LoginUser(payload)

	if err1 != nil {
		utils.SendMsg(err1.(*utils.HttpError).Message, err1.(*utils.HttpError).Status, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{Msg: "Tokens generated", Data: tokens}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      User Register
// @Description  API for registering users
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body service.RegisterUserPayload true "Register user payload"
// @Success      201  {object} utils.ResUser
// @Failure      400  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /register [post]
func RegisterController(w http.ResponseWriter, r *http.Request) {
	var payload service.RegisterUserPayload

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		fmt.Println("Invalid payload", err.Error())
		utils.SendMsg("Invalid payload", http.StatusBadRequest, w)
		return
	}

	valErr := utils.Validate(payload)

	if valErr != nil {
		fmt.Println("Invalid payload", valErr.Error())
		utils.SendMsg("Invalid payload", http.StatusBadRequest, w)
		return
	}

	str, err1 := service.RegisterUser(payload)

	if err1 != nil {
		utils.SendMsg(err1.(*utils.HttpError).Message, err1.(*utils.HttpError).Status, w)
		return
	}

	utils.SendMsg(str, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Renew Access Token
// @Description  API for renewing the access token. Make sure to provide the refresh token to get the new tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        token body service.RenewAccessTokenPayload true "Renew access token payload"
// @Success      201  {object} utils.ResUserWithData{data=lib.Tokens}
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /renew-access-token [get]
func RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	var token service.RenewAccessTokenPayload

	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		fmt.Println("Invalid payload", err.Error())
		utils.SendMsg("Invalid payload", http.StatusBadRequest, w)
		return
	}

	newTokens, err1 := service.RenewAccessToken(token)

	if err1 != nil {
		utils.SendMsg(err1.(*utils.HttpError).Message, err1.(*utils.HttpError).Status, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{Msg: "Tokens generated", Data: newTokens}, http.StatusCreated, w)
}

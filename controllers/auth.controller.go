package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	service "github.com/kennizen/e-commerce-backend/services"
	"github.com/kennizen/e-commerce-backend/utils"
)

// @Summary      User Login
// @Description  API for user login using email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body service.LoginUserPayload true "Login credentials"
// @Success      200  {object} utils.ResUserWithData
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

	service.LoginUser(payload, w)
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

	service.RegisterUser(payload, w)
}

// ---------------------------------------------------------------------------------------- //

// @Summary      Renew Access Token
// @Description  API for renewing the access token. Make sure to provide the refresh token to get the new tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer refreshToken"
// @Success      201  {object} utils.ResUserWithData
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /renew-access-token [get]
func RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if token == "" {
		utils.SendMsg("Invalid token", http.StatusUnauthorized, w)
		return
	}

	bearerToken := strings.Split(token, " ")

	if bearerToken[0] != "Bearer" && bearerToken[1] == "" {
		utils.SendMsg("Invalid token", http.StatusUnauthorized, w)
		return
	}

	service.RenewAccessToken(bearerToken[1], w)
}

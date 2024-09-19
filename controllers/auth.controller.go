package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	service "github.com/kennizen/e-commerce-backend/services"
	"github.com/kennizen/e-commerce-backend/utils"
)

type AccessToken struct {
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}

// @Summary      User Login
// @Description  API for user login using email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body service.LoginUserPayload true "Login credentials"
// @Success      200  {object} utils.ResUserWithData{data=AccessToken}
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

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/renew-access-token",
		HttpOnly: true,
		Expires:  time.Unix(tokens.RefreshTokenExp, 0),
	})

	utils.SendJson(utils.ResUserWithData{Msg: "Token", Data: AccessToken{
		Token:  tokens.Token,
		Expiry: tokens.TokenExp,
	}}, http.StatusCreated, w)
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
// @Description  API for renewing the access token. The refresh token is set in the http-only cookie when the user first logs in so renewing the token will only work if the user have logged in atleast once.
// @Tags         Auth
// @Produce      json
// @Success      201  {object} utils.ResUserWithData{data=AccessToken}
// @Failure      401  {object} utils.ResUser
// @Failure      500  {object} utils.ResUser
// @Router       /renew-access-token [get]
func RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		utils.SendMsg("Invalid token", http.StatusUnauthorized, w)
		return
	}

	newTokens, err1 := service.RenewAccessToken(cookie.Value)

	if err1 != nil {
		utils.SendMsg(err1.(*utils.HttpError).Message, err1.(*utils.HttpError).Status, w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newTokens.RefreshToken,
		Path:     "/renew-access-token",
		HttpOnly: true,
		Expires:  time.Unix(newTokens.RefreshTokenExp, 0),
	})

	utils.SendJson(utils.ResUserWithData{Msg: "Token", Data: AccessToken{
		Token:  newTokens.Token,
		Expiry: newTokens.TokenExp,
	}}, http.StatusCreated, w)
}

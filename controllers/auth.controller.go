package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	service "github.com/kennizen/e-commerce-backend/services"
	"github.com/kennizen/e-commerce-backend/utils"
)

// @Summary      List accounts
// @Description  get accounts
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Failure      400
// @Failure      404
// @Failure      500
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

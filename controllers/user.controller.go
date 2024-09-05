package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	var user service.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.SendMsg(err.Error(), http.StatusInternalServerError, w)
		return
	}

	fmt.Println("user", user)

	user.LoginUser(w)
}

// ---------------------------------------------------------------------------------------- //

func RegisterController(w http.ResponseWriter, r *http.Request) {
	var user service.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.SendMsg(err.Error(), http.StatusInternalServerError, w)
		return
	}

	fmt.Println("user", user)

	user.RegisterUser(w)
}

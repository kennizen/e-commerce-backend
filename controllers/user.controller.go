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
	var payload service.LoginUserPayload

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		utils.SendMsg(err.Error(), http.StatusInternalServerError, w)
		return
	}

	fmt.Println("user", payload)

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

	fmt.Println("user", payload)

	service.RegisterUser(payload, w)
}

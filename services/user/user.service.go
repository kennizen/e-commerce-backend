package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/utils"
)

type User struct {
	Firstname  string `json:"firstname"`
	Middlename string `json:"middlename"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
	Age        int    `json:"age"`
	Avatar     string `json:"avatar"`
	Password   string `json:"password"`
}

func (u *User) RegisterUserService(w http.ResponseWriter) {
	var hashed_password string

	h := sha256.New()
	h.Write([]byte(u.Password))

	hashed_password = hex.EncodeToString(h.Sum(nil))

	// check if user already exists
	rows, err := db.DB.Query("SELECT id FROM customers WHERE email = $1", u.Email)

	if err != nil {
		fmt.Println("Failed to query customers", err.Error())
		utils.SendResp("Server error", http.StatusInternalServerError, w)
		return
	}

	if data := rows.Next(); data {
		utils.SendResp("User already exists", http.StatusConflict, w)
		return
	}

	// if not user already present insert the user
	trx, err := db.DB.Begin()

	if err != nil {
		fmt.Println("Failed to start a trx", err.Error())
		utils.SendResp("Server error", http.StatusInternalServerError, w)
		return
	}

	var newAge int

	if u.Age <= 0 {
		newAge = 1
	} else {
		newAge = u.Age
	}
	_, err1 := trx.Exec("INSERT INTO customers (firstname, middlename, lastname, email, age, avatar, hashed_password) VALUES ($1, $2, $3, $4, $5, $6, $7)", u.Firstname, u.Middlename, u.Lastname, u.Email, newAge, u.Avatar, hashed_password)

	if err1 != nil {
		fmt.Println("Failed to insert in customers", err1.Error())
		trx.Rollback()
		utils.SendResp("Server error", http.StatusInternalServerError, w)
		return
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Failed to commit", comErr.Error())
		utils.SendResp("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendResp("User created", http.StatusCreated, w)
}

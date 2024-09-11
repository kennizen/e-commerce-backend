package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/lib"
	"github.com/kennizen/e-commerce-backend/utils"
)

type RegisterUserPayload struct {
	Firstname  string `validate:"required"`
	Middlename string
	Lastname   string `validate:"required"`
	Age        int    `validate:"required,gte=1"`
	Email      string `validate:"required,email"`
	Password   string `validate:"required"`
}

type LoginUserPayload struct {
	Email    string
	Password string
}

func RegisterUser(arg RegisterUserPayload, w http.ResponseWriter) {
	var hashed_password string

	h := sha256.New()
	h.Write([]byte(arg.Password))

	hashed_password = hex.EncodeToString(h.Sum(nil))

	// check if user already exists
	rows, err := db.DB.Query("SELECT id FROM customers WHERE email = $1", arg.Email)

	if err != nil {
		fmt.Println("Failed to query customers", err.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	defer rows.Close()

	if rows.Next() {
		utils.SendMsg("User already exists", http.StatusConflict, w)
		return
	}

	// if not user already present insert the user
	trx, err := db.DB.Begin()

	if err != nil {
		fmt.Println("Failed to start a trx", err.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	_, err1 := trx.Exec(
		"INSERT INTO customers (firstname, middlename, lastname, email, age, hashed_password) VALUES ($1, $2, $3, $4, $5, $6)", arg.Firstname, arg.Middlename, arg.Lastname, arg.Email, arg.Age, hashed_password,
	)

	if err1 != nil {
		fmt.Println("Failed to insert in customers", err1.Error())
		trx.Rollback()
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Failed to commit", comErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendMsg("User created", http.StatusCreated, w)
}

// ---------------------------------------------------------------------------------------- //

func LoginUser(arg LoginUserPayload, w http.ResponseWriter) {
	// check if user exists and correct password
	var hashed_password string

	h := sha256.New()
	h.Write([]byte(arg.Password))

	hashed_password = hex.EncodeToString(h.Sum(nil))

	rows, err := db.DB.Query("SELECT id, email FROM customers WHERE email = $1 and hashed_password = $2", arg.Email, hashed_password)

	if err != nil {
		fmt.Println("Failed to query customers", err.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	defer rows.Close()

	isEmpty := true
	var id int
	var email string

	for rows.Next() {
		isEmpty = false
		err := rows.Scan(&id, &email)
		if err != nil {
			log.Fatalln("Error scanning row", err.Error())
		}
	}

	if isEmpty {
		utils.SendMsg("User not found", http.StatusNotFound, w)
		return
	}

	// all success then send access and refresh token
	tokens, err1 := lib.GenerateTokens(strconv.Itoa(id), email)

	if err1 != nil {
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(tokens, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func RenewAccessToken(refToken string, w http.ResponseWriter) {
	claims, isValid := lib.ValidateToken(refToken, os.Getenv("JWT_REFRESH_TOKEN_SECRET"))

	if !isValid {
		utils.SendMsg("Invalid token", http.StatusUnauthorized, w)
		return
	}

	newTokens, err := lib.GenerateTokens(claims.Id, claims.Email)

	if err != nil {
		fmt.Println("error in generating new tokens")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(newTokens, http.StatusCreated, w)
}

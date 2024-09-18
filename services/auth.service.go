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
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

type RenewAccessTokenPayload struct {
	RefToken string `validate:"required"`
}

func RegisterUser(arg RegisterUserPayload) (string, error) {
	var hashed_password string

	h := sha256.New()
	h.Write([]byte(arg.Password))

	hashed_password = hex.EncodeToString(h.Sum(nil))

	// check if user already exists
	rows, err := db.DB.Query("SELECT id FROM customers WHERE email = $1", arg.Email)

	if err != nil {
		fmt.Println("Failed to query customers", err.Error())
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	defer rows.Close()

	if rows.Next() {
		return "", utils.NewHttpError("User already exists", http.StatusConflict)
	}

	// if not user already present insert the user
	trx, err := db.DB.Begin()

	if err != nil {
		fmt.Println("Failed to start a trx", err.Error())
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	_, err1 := trx.Exec(
		"INSERT INTO customers (firstname, middlename, lastname, email, age, hashed_password) VALUES ($1, $2, $3, $4, $5, $6)", arg.Firstname, arg.Middlename, arg.Lastname, arg.Email, arg.Age, hashed_password,
	)

	if err1 != nil {
		fmt.Println("Failed to insert in customers", err1.Error())
		trx.Rollback()
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Failed to commit", comErr.Error())
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	return "User created", nil
}

// ---------------------------------------------------------------------------------------- //

func LoginUser(payload LoginUserPayload) (*lib.Tokens, error) {
	// check if user exists and correct password
	var hashed_password string

	h := sha256.New()
	h.Write([]byte(payload.Password))

	hashed_password = hex.EncodeToString(h.Sum(nil))

	rows, err := db.DB.Query("SELECT id, email FROM customers WHERE email = $1 and hashed_password = $2", payload.Email, hashed_password)

	if err != nil {
		fmt.Println("Failed to query customers", err.Error())
		return &lib.Tokens{}, utils.NewHttpError("Server error", http.StatusInternalServerError)
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
		return &lib.Tokens{}, utils.NewHttpError("User not found", http.StatusNotFound)
	}

	// all success then send access and refresh token
	tokens, err1 := lib.GenerateTokens(strconv.Itoa(id), email)

	if err1 != nil {
		return &lib.Tokens{}, utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	claims, _ := lib.ValidateToken(tokens.RefreshToken, os.Getenv("JWT_REFRESH_TOKEN_SECRET"))

	upErr := updateRefTokenId(claims.ID, claims.Email)

	if upErr != nil {
		fmt.Println("Error in transaction", upErr.Error())
		return &lib.Tokens{}, utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	return &tokens, nil
}

// ---------------------------------------------------------------------------------------- //

func RenewAccessToken(refToken RenewAccessTokenPayload) (*lib.Tokens, error) {
	claims, isValid := lib.ValidateToken(refToken.RefToken, os.Getenv("JWT_REFRESH_TOKEN_SECRET"))

	if !isValid {
		return &lib.Tokens{}, utils.NewHttpError("Invalid token", http.StatusUnauthorized)
	}

	var refTokenId string

	row := db.DB.QueryRow("SELECT reftoken_id FROM customers WHERE email = $1", claims.Email)
	row.Scan(&refTokenId)

	if refTokenId == "" || refTokenId != claims.ID {
		fmt.Println("Invalid refresh token")
		return &lib.Tokens{}, utils.NewHttpError("Invalid refresh token", http.StatusUnauthorized)
	}

	newTokens, err := lib.GenerateTokens(claims.Id, claims.Email)

	newClaims, _ := lib.ValidateToken(newTokens.RefreshToken, os.Getenv("JWT_REFRESH_TOKEN_SECRET"))

	upErr := updateRefTokenId(newClaims.ID, newClaims.Email)

	if upErr != nil {
		fmt.Println("Error in transaction", upErr.Error())
		return &lib.Tokens{}, utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	if err != nil {
		fmt.Println("error in generating new tokens")
		return &lib.Tokens{}, utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	return &newTokens, nil
}

// ---------------------------------------------------------------------------------------- //

func updateRefTokenId(id, email string) error {
	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		return utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	_, execErr := trx.Exec("UPDATE customers SET reftoken_id = $1 WHERE email = $2", id, email)

	if execErr != nil {
		fmt.Println("Failed to update customers")
		return utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		return utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	return nil
}

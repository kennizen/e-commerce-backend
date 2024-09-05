package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/lib"
	"github.com/kennizen/e-commerce-backend/utils"
)

type User struct {
	Id         int    `json:"-"`
	Firstname  string `json:"firstname"`
	Middlename string `json:"middlename"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
	Age        int    `json:"age"`
	Avatar     string `json:"avatar"`
	Password   string `json:"password"`
	Created_at string `json:"-"`
	Updated_at string `json:"-"`
}

func (u *User) RegisterUser(w http.ResponseWriter) {
	var hashed_password string

	h := sha256.New()
	h.Write([]byte(u.Password))

	hashed_password = hex.EncodeToString(h.Sum(nil))

	// check if user already exists
	rows, err := db.DB.Query("SELECT id FROM customers WHERE email = $1", u.Email)

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

func (u *User) LoginUser(w http.ResponseWriter) {
	// check if user exists and correct password
	var hashed_password string

	h := sha256.New()
	h.Write([]byte(u.Password))

	hashed_password = hex.EncodeToString(h.Sum(nil))

	rows, err := db.DB.Query("SELECT id, email FROM customers WHERE email = $1 and hashed_password = $2", u.Email, hashed_password)

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

	fmt.Println("id and email", id, email)

	// all success then send access and refresh token
	tokens, err1 := lib.GenerateTokens(strconv.Itoa(id), email)

	if err1 != nil {
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(tokens, http.StatusOK, w)
}

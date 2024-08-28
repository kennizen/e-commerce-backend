package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func (u *User) RegisterUserService() {
	var hashed_password string

	h := sha256.New()
	h.Write([]byte(u.Password))

	hashed_password = hex.EncodeToString(h.Sum(nil))

	// check if user already exists
	fmt.Println("pass", hashed_password)

}

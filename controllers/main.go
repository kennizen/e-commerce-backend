package controllers

import (
	"fmt"
	"net/http"
)

func LoginController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("the handler in working")
}

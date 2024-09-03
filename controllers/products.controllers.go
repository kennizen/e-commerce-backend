package controller

import (
	"fmt"
	"net/http"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting all products")
}

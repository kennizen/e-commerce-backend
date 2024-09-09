package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/middlewares"
	"github.com/kennizen/e-commerce-backend/routes"
)

func main() {
	db.InitDB()
	router := http.NewServeMux()
	routes.RegisterRoutes(router)

	defer db.DB.Close()

	server := http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: middlewares.Logger(router),
	}

	fmt.Println("🚀🚀🚀 Server started on host port", os.Getenv("API_PORT"))
	log.Fatal(server.ListenAndServe())
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/routes"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		fmt.Println("Successfully loaded .env file.")
	}
}

func main() {
	// loadEnv()
	db.InitDB()
	router := http.NewServeMux()
	routes.RegisterRoutes(router)

	fmt.Println("🚀🚀🚀 Server started on port", os.Getenv("API_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("API_PORT"), router))
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/routes"
)

func main() {
	db.InitDB()
	router := http.NewServeMux()
	routes.RegisterRoutes(router)

	server := http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: router,
	}

	fmt.Println("🚀🚀🚀 Server started on port", os.Getenv("SERVER_PORT"))
	log.Fatal(server.ListenAndServe())
}

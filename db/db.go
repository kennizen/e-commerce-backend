package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func createDB() *sql.DB {
	connStr := "host=postgres port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error connecting database")
	}

	// check to see if ecommerce db already exists
	var exists bool
	query := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = 'ecommerce')`
	err = db.QueryRow(query).Scan(&exists)

	if err != nil {
		log.Fatal(err)
	}

	// If the database doesn't exist, create it
	if !exists {
		_, err := db.Exec("CREATE DATABASE ecommerce")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Database 'ecommerce' created successfully.")
	} else {
		fmt.Println("Database 'ecommerce' already exists.")
	}

	db.Close()

	ecomConnStr := "host=postgres port=5432 user=postgres password=postgres dbname=ecommerce sslmode=disable"
	ecommerceDB, err := sql.Open("postgres", ecomConnStr)

	fmt.Println(ecomConnStr)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to ecommerce db.")
	}

	return ecommerceDB
}

func runMigrations() {
	m, err := migrate.New("file://db/migrations", "postgres://postgres:postgres@postgres:5432/ecommerce?sslmode=disable")

	if err != nil {
		log.Fatal("Error in connecting db during migration", err)
	}

	if err := m.Up(); err != nil {
		fmt.Println("Migrations", err)
	} else {
		fmt.Println("Migrations runned successfully.")
	}

	m.Close()
}

func InitDB() {
	DB = createDB()
	runMigrations()
}

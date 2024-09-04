package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/utils"
)

type Product struct {
	Id           int     `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	Price        float32 `json:"price"`
	Stock        int     `json:"stock"`
	Image        string  `json:"image"`
	Thumbnail    string  `json:"thumbnail"`
	Rating       float32 `json:"rating"`
	Weight       int     `json:"weight"`
	Width        float32 `json:"width"`
	Height       float32 `json:"height"`
	Depth        float32 `json:"depth"`
	Warranty     string  `json:"warranty"`
	Shipping     string  `json:"shipping"`
	Availability string  `json:"availability"`
	ReturnPolicy string  `json:"returnPolicy"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

func GetProducts(page, limit int, w http.ResponseWriter) {
	rows, err := db.DB.Query("SELECT * FROM products OFFSET ($1 - 1) * $2 LIMIT $2", page, limit)
	rowCount := db.DB.QueryRow("SELECT COUNT(id) FROM products")

	if err != nil {
		fmt.Println("Failed to query customers", err.Error())
		utils.SendError("Server error", http.StatusInternalServerError, w)
		return
	}

	defer rows.Close()

	var count int
	rowCount.Scan(&count)
	product := Product{}
	products := make([]Product, 0)
	isEmpty := true

	for rows.Next() {
		isEmpty = false
		err := rows.Scan(
			&product.Id,
			&product.Title,
			&product.Description,
			&product.Category,
			&product.Price,
			&product.Stock,
			&product.Image,
			&product.Thumbnail,
			&product.Rating,
			&product.Weight,
			&product.Width,
			&product.Height,
			&product.Depth,
			&product.Warranty,
			&product.Shipping,
			&product.Availability,
			&product.ReturnPolicy,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			log.Fatalln("Error scanning row", err.Error())
		}

		products = append(products, product)
	}

	if isEmpty {
		utils.SendJson(map[string][]any{"data": make([]any, 0)}, http.StatusOK, w)
		return
	}

	utils.SendJson(map[string]any{"data": products, "totalCount": count}, http.StatusOK, w)
}

func GetProduct(id int, w http.ResponseWriter) {
	row := db.DB.QueryRow("SELECT * FROM products WHERE id = $1", id)

	if row == nil {
		utils.SendError("Product not found", http.StatusNotFound, w)
		return
	}

	product := Product{}

	err := row.Scan(
		&product.Id,
		&product.Title,
		&product.Description,
		&product.Category,
		&product.Price,
		&product.Stock,
		&product.Image,
		&product.Thumbnail,
		&product.Rating,
		&product.Weight,
		&product.Width,
		&product.Height,
		&product.Depth,
		&product.Warranty,
		&product.Shipping,
		&product.Availability,
		&product.ReturnPolicy,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		log.Fatalln("Error scanning row", err.Error())
	}

	utils.SendJson(map[string]any{"data": product}, http.StatusOK, w)
}

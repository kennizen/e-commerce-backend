package service

import (
	"fmt"
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

	defer rows.Close()

	if err != nil {
		fmt.Println("Failed to query customers", err.Error())
		utils.SendError("Server error", http.StatusInternalServerError, w)
		return
	}

	// product := Product{}

	for rows.Next() {
		// err := rows.Scan(&product.Id, &product.Description)
	}
}

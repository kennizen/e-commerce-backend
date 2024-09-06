package service

import (
	"database/sql"
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
	CreatedAt    string  `json:"-"`
	UpdatedAt    string  `json:"-"`
}

type ProductReview struct {
	Id        int
	ReviewBy  int
	ProductId int
	Review    string
	Rating    float32
	CreatedAt string
}

type ProductReviewArgs struct {
	Review    string
	Rating    float32
	ProductId string
	UserId    string
}

// ---------------------------------------------------------------------------------------- //

func GetProducts(page, limit int, w http.ResponseWriter) {
	var rows *sql.Rows
	var err error

	if limit == -1 {
		rows, err = db.DB.Query("SELECT * FROM products")
	} else {
		rows, err = db.DB.Query("SELECT * FROM products OFFSET ($1 - 1) * $2 LIMIT $2", page, limit)
	}

	rowCount := db.DB.QueryRow("SELECT COUNT(id) FROM products")

	if err != nil {
		fmt.Println("Failed to query customers", err.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
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

// ---------------------------------------------------------------------------------------- //

func GetProduct(id int, w http.ResponseWriter) {
	row := db.DB.QueryRow("SELECT * FROM products WHERE id = $1", id)

	if row == nil {
		utils.SendMsg("Product not found", http.StatusNotFound, w)
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

// ---------------------------------------------------------------------------------------- //

func MarkFavorite(userId, productId string, w http.ResponseWriter) {
	rows, err := db.DB.Query("SELECT id FROM favorites WHERE customer_id = $1 AND product_id = $2", userId, productId)

	if rows.Next() {
		utils.SendMsg("Product already added to favorites", http.StatusConflict, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	_, err1 := trx.Exec("INSERT INTO favorites (customer_id, product_id) VALUES ($1, $2)", userId, productId)

	if err1 != nil || err != nil {
		fmt.Println("Error inserting data in favorite", err1.Error(), err.Error())
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	rows.Close()

	utils.SendJson(map[string]string{"data": "product added to favorites"}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func UnMarkFavorite(userId, productId string, w http.ResponseWriter) {
	var favId string = ""

	row := db.DB.QueryRow("SELECT id FROM favorites WHERE customer_id = $1 AND product_id = $2", userId, productId)
	row.Scan(&favId)

	if favId == "" {
		fmt.Println("Product not found to unmark favorite.")
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	_, err := trx.Exec("DELETE FROM favorites WHERE id = $1", favId)

	if err != nil {
		fmt.Println("Failed to unmark favorite product.")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(map[string]string{"data": "product removed from favorites"}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func GetFavorites(userId string, w http.ResponseWriter) {
	rows, err := db.DB.Query(
		`select 
				p.id, 
				p.title, 
				p.description, 
				p.category, 
				p.price, 
				p.stock, 
				p.image, 
				p.thumbnail, 
				p.rating, 
				p.weight, 
				p.width, 
				p.height, 
				p."depth", 
				p.warranty, 
				p.shipping, 
				p.availability,
				p.return_policy 
			from products p 
			left join favorites f on f.product_id  = p.id
			where f.customer_id = $1`,
		userId,
	)

	if err != nil {
		fmt.Println("Failed to get favorites")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	var products []Product = make([]Product, 0)
	var product Product

	for rows.Next() {
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
		)

		if err != nil {
			log.Fatalln("Error scanning row", err.Error())
		}

		products = append(products, product)
	}

	utils.SendJson(map[string][]Product{"data": products}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func AddProductReview(args ProductReviewArgs, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow(
		"SELECT id FROM product_reviews WHERE review_by = $1 AND product_id = $2", args.UserId, args.ProductId,
	)
	row.Scan(&id)

	if id != "" {
		fmt.Println("Review already added by user", args.UserId)
		utils.SendMsg("Already reviewed the product", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	_, err := trx.Exec(
		"INSERT INTO product_reviews (review_by, product_id, review, rating) VALUES ($1, $2, $3, $4)", args.UserId, args.ProductId, args.Review, args.Rating,
	)

	if err != nil {
		fmt.Println("Failed to add product review")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendMsg("Review Added", http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func UpdateProductReview(review ProductReview, w http.ResponseWriter) {

}

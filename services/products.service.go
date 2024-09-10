package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/models"
	"github.com/kennizen/e-commerce-backend/utils"
)

type ProductReviewArgs struct {
	Review    string
	Rating    float32
	ProductId string
	UserId    string
}

type ProductUpdateArgs struct {
	ReviewId string
	Review   string
	Rating   float32
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
	product := models.Product{}
	products := make([]models.Product, 0)
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

	product := models.Product{}

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

	utils.SendJson(map[string]string{"data": "Product added to favorites"}, http.StatusOK, w)
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

	utils.SendJson(map[string]string{"data": "Product removed from favorites"}, http.StatusOK, w)
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

	var products []models.Product = make([]models.Product, 0)
	var product models.Product

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

	utils.SendJson(map[string][]models.Product{"data": products}, http.StatusOK, w)
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

	row1 := trx.QueryRow(
		"INSERT INTO product_reviews (review_by, product_id, review, rating) VALUES ($1, $2, $3, $4) RETURNING *",
		args.UserId, args.ProductId, args.Review, args.Rating,
	)

	var proReview models.ProductReview

	row1.Scan(
		&proReview.Id,
		&proReview.ReviewBy,
		&proReview.ProductId,
		&proReview.Review,
		&proReview.Rating,
		&proReview.CreatedAt,
		&proReview.UpdatedAt,
	)

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(map[string]interface{}{
		"message": "Product review added",
		"data":    proReview,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func UpdateProductReview(args ProductUpdateArgs, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM product_reviews WHERE id = $1", args.ReviewId)
	row.Scan(&id)

	if id == "" {
		fmt.Println("Review not found to update", args.ReviewId)
		utils.SendMsg("Review not found to update", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	row1 := trx.QueryRow(
		"UPDATE product_reviews SET review = $1, rating = $2 WHERE id = $3 RETURNING *",
		args.Review, args.Rating, args.ReviewId,
	)

	var proReview models.ProductReview

	row1.Scan(
		&proReview.Id,
		&proReview.ReviewBy,
		&proReview.ProductId,
		&proReview.Review,
		&proReview.Rating,
		&proReview.CreatedAt,
		&proReview.UpdatedAt,
	)

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(map[string]interface{}{
		"message": "Product review updated",
		"data":    proReview,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func DeleteProductReview(reviewId string, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM product_reviews WHERE id = $1", reviewId)
	row.Scan(&id)

	if id == "" {
		fmt.Println("Review not found to delete", reviewId)
		utils.SendMsg("Review not found to delete", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	delRow := trx.QueryRow("DELETE FROM product_reviews WHERE id = $1 RETURNING *", reviewId)

	var proReview models.ProductReview

	delRow.Scan(
		&proReview.Id,
		&proReview.ReviewBy,
		&proReview.ProductId,
		&proReview.Review,
		&proReview.Rating,
		&proReview.CreatedAt,
		&proReview.UpdatedAt,
	)

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(map[string]interface{}{
		"message": "Product review deleted",
		"data":    proReview,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func GetProductReviewsByProductId(productId string, w http.ResponseWriter) {
	rows, err := db.DB.Query(
		`select 
			pr.id, 
			pr.review,
			pr.rating, 
			pr.created_at, 
			pr.updated_at, 
			c.id as customer_id, 
			c.firstname, 
			c.middlename, 
			c.lastname, 
			c.email, 
			c.age, 
			c.avatar 
		from product_reviews pr 
		right join customers c on pr.review_by = c.id 
		where pr.product_id = $1`,
		productId,
	)

	if err != nil {
		fmt.Println("Failed to execute query")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	type Review struct {
		Id        string
		Review    string
		Rating    float32
		CreatedAt string
		UpdatedAt string
	}

	type Customer struct {
		Id         string
		Firstname  string
		Middlename string
		Lastname   string
		Email      string
		Age        string
		Avatar     string
	}

	type Data struct {
		Review   Review
		Customer Customer
	}

	var resp []Data
	isEmpty := true

	for rows.Next() {
		isEmpty = false
		var rev = Review{}
		var cust = Customer{}

		err := rows.Scan(
			&rev.Id,
			&rev.Review,
			&rev.Rating,
			&rev.CreatedAt,
			&rev.UpdatedAt,
			&cust.Id,
			&cust.Firstname,
			&cust.Middlename,
			&cust.Lastname,
			&cust.Email,
			&cust.Age,
			&cust.Avatar,
		)

		if err != nil {
			log.Fatalln("Error scanning row", err.Error())
		}

		resp = append(resp, Data{
			Review:   rev,
			Customer: cust,
		})
	}

	if isEmpty {
		utils.SendJson(map[string][]any{"data": make([]any, 0)}, http.StatusOK, w)
		return
	}

	utils.SendJson(resp, http.StatusOK, w)
}

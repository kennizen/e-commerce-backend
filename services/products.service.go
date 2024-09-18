package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/models"
	"github.com/kennizen/e-commerce-backend/utils"
)

type ProductReviewPayload struct {
	Review string  `validate:"required"`
	Rating float32 `validate:"required"`
}

type ProductsResponse struct {
	Products   *[]models.Product
	TotalCount int
}

// ---------------------------------------------------------------------------------------- //

func GetProducts(page, limit int) (*ProductsResponse, error) {
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
		return nil, utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	defer rows.Close()

	var count int = 0
	rowCount.Scan(&count)
	product := models.Product{}
	products := make([]models.Product, 0)

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
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			log.Fatalln("Error scanning row", err.Error())
		}

		products = append(products, product)
	}

	return &ProductsResponse{
		Products:   &products,
		TotalCount: count,
	}, nil
}

// ---------------------------------------------------------------------------------------- //

func GetProduct(id int) (*models.Product, error) {
	row := db.DB.QueryRow("SELECT * FROM products WHERE id = $1", id)

	if row == nil {
		return nil, utils.NewHttpError("Product not found", http.StatusNotFound)
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

	return &product, nil
}

// ---------------------------------------------------------------------------------------- //

func MarkFavorite(userId, productId string) (string, error) {
	rows, err := db.DB.Query("SELECT id FROM favorites WHERE customer_id = $1 AND product_id = $2", userId, productId)

	if rows.Next() {
		return "", utils.NewHttpError("Product already added to favorites", http.StatusConflict)
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	_, err1 := trx.Exec("INSERT INTO favorites (customer_id, product_id) VALUES ($1, $2)", userId, productId)

	if err1 != nil || err != nil {
		fmt.Println("Error inserting data in favorite", err1.Error(), err.Error())
		return "", utils.NewHttpError("Bad request", http.StatusBadRequest)
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	rows.Close()

	return "Product added to favorites", nil
}

// ---------------------------------------------------------------------------------------- //

func UnMarkFavorite(userId, productId string) (string, error) {
	var favId string = ""

	row := db.DB.QueryRow("SELECT id FROM favorites WHERE customer_id = $1 AND product_id = $2", userId, productId)
	row.Scan(&favId)

	if favId == "" {
		fmt.Println("Product not found to unmark favorite.")
		return "", utils.NewHttpError("Bad request", http.StatusBadRequest)
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	_, err := trx.Exec("DELETE FROM favorites WHERE id = $1", favId)

	if err != nil {
		fmt.Println("Failed to unmark favorite product.")
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	return "Product removed from favorites", nil
}

// ---------------------------------------------------------------------------------------- //

func GetFavorites(userId string) (*[]models.Product, error) {
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
		return nil, utils.NewHttpError("Server error", http.StatusInternalServerError)
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

	return &products, nil
}

// ---------------------------------------------------------------------------------------- //

func AddProductReview(args ProductReviewPayload, userId, productId string) (*models.ProductReview, error) {
	var id string = ""

	row := db.DB.QueryRow(
		"SELECT id FROM product_reviews WHERE review_by = $1 AND product_id = $2", userId, productId,
	)
	row.Scan(&id)

	if id != "" {
		fmt.Println("Review already added by user", userId)
		return nil, utils.NewHttpError("Already reviewed the product", http.StatusBadRequest)
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		return nil, utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	row1 := trx.QueryRow(
		"INSERT INTO product_reviews (review_by, product_id, review, rating) VALUES ($1, $2, $3, $4) RETURNING *",
		userId, productId, args.Review, args.Rating,
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
		return nil, utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	return &proReview, nil
}

// ---------------------------------------------------------------------------------------- //

func UpdateProductReview(args ProductReviewPayload, userId, reviewId string, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM product_reviews WHERE id = $1", reviewId)
	row.Scan(&id)

	if id == "" {
		fmt.Println("Review not found to update", reviewId)
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
		"UPDATE product_reviews SET review = $1, rating = $2, updated_at = $3 WHERE id = $4 AND review_by = $5 RETURNING *",
		args.Review, args.Rating, time.Now().UTC().Format(time.RFC3339), reviewId, userId,
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

	utils.SendJson(utils.ResUserWithData{
		Msg:  "Product review updated",
		Data: proReview,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func DeleteProductReview(reviewId, userId string, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM product_reviews WHERE id = $1 AND review_by = $2", reviewId, userId)
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

	delRow := trx.QueryRow("DELETE FROM product_reviews WHERE id = $1 AND review_by = $2 RETURNING *", reviewId, userId)

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

	utils.SendJson(utils.ResUserWithData{
		Msg:  "Product review deleted",
		Data: proReview,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

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

type AllReviewsResponse struct {
	Review   Review
	Customer Customer
}

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

	var resp []AllReviewsResponse
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

		resp = append(resp, AllReviewsResponse{
			Review:   rev,
			Customer: cust,
		})
	}

	if isEmpty {
		utils.SendJson(utils.ResUserWithData{
			Msg:  "No reviews found",
			Data: make([]any, 0),
		}, http.StatusOK, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{
		Msg:  "No reviews found",
		Data: resp,
	}, http.StatusOK, w)
}

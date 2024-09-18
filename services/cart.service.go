package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/models"
	"github.com/kennizen/e-commerce-backend/utils"
)

type AddToCartPayload struct {
	Quantity int `validate:"required,gte="`
}

func AddToCart(args AddToCartPayload, userId, productId string) (string, error) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM cart WHERE customer_id = $1 AND product_id = $2", userId, productId)
	row.Scan(&id)

	if id != "" {
		fmt.Println("Product already added to cart")
		return "", utils.NewHttpError("Product already in cart", http.StatusBadRequest)
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	_, err := trx.Exec(
		"INSERT INTO cart (customer_id, product_id, quantity) VALUES ($1, $2, $3)",
		userId, productId, args.Quantity,
	)

	if err != nil {
		fmt.Println("Failed to add product to cart.")
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	return "Product added to cart", nil
}

// ---------------------------------------------------------------------------------------- //

func RemoveFromCart(userId, productId string) (string, error) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM cart WHERE customer_id = $1 AND product_id = $2", userId, productId)
	row.Scan(&id)

	if id == "" {
		fmt.Println("Product not found to remove")
		return "", utils.NewHttpError("Product not found to remove", http.StatusBadRequest)
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	_, err := trx.Exec("DELETE FROM cart WHERE id = $1", id)

	if err != nil {
		fmt.Println("Failed to remove product from cart.")
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	return "Product removed from cart", nil
}

// ---------------------------------------------------------------------------------------- //

func UpdateCartItems(args AddToCartPayload, userId, productId string) (string, error) {
	var id string

	row := db.DB.QueryRow(
		"SELECT id FROM cart WHERE customer_id = $1 AND product_id = $2", userId, productId,
	)
	row.Scan(&id)

	if id == "" {
		fmt.Println("Product not found to update")
		return "", utils.NewHttpError("Product not found to update", http.StatusBadRequest)
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	_, err := trx.Exec(
		"UPDATE cart SET quantity = $1, updated_at = $2 WHERE id = $3",
		args.Quantity, time.Now().UTC().Format(time.RFC3339), id,
	)

	if err != nil {
		fmt.Println("Failed to update cart.")
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		return "", utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	return "Cart updated", nil
}

// ---------------------------------------------------------------------------------------- //

type GetCartResponse struct {
	models.Product
	Quantity int
}

func GetCart(userId string) (*[]GetCartResponse, error) {
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
				p.return_policy,
				c.quantity 
			from products p 
			left join cart c on c.product_id  = p.id
			where c.customer_id = $1`,
		userId,
	)

	if err != nil {
		fmt.Println("Failed to get cart")
		return nil, utils.NewHttpError("Server error", http.StatusInternalServerError)
	}

	var res []GetCartResponse = make([]GetCartResponse, 0)
	var product GetCartResponse

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
			&product.Quantity,
		)

		if err != nil {
			log.Fatalln("Error scanning row", err.Error())
		}

		res = append(res, product)
	}

	return &res, nil
}

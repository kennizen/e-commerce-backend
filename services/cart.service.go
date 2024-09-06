package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/utils"
)

type CartArgs struct {
	UserId    string
	ProductId string
	Quantity  int
}

func AddToCart(args CartArgs, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM cart WHERE customer_id = $1 AND product_id = $2", args.UserId, args.ProductId)
	row.Scan(&id)

	if id != "" {
		fmt.Println("Product already added to cart")
		utils.SendMsg("Product already in cart", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	_, err := trx.Exec(
		"INSERT INTO cart (customer_id, product_id, quantity) VALUES ($1, $2, $3)",
		args.UserId, args.ProductId, args.Quantity,
	)

	if err != nil {
		fmt.Println("Failed to add product to cart.")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendMsg("Product added to cart", http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func RemoveFromCart(userId, productId string, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM cart WHERE customer_id = $1 AND product_id = $2", userId, productId)
	row.Scan(&id)

	if id == "" {
		fmt.Println("Product not found to remove")
		utils.SendMsg("Product not found to remove", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	_, err := trx.Exec("DELETE FROM cart WHERE id = $1", id)

	if err != nil {
		fmt.Println("Failed to remove product from cart.")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendMsg("Product removed from cart", http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func UpdateCartItems(args CartArgs, w http.ResponseWriter) {
	var result struct {
		Id       string
		Quantity int
	}

	row := db.DB.QueryRow(
		"SELECT id, quantity FROM cart WHERE customer_id = $1 AND product_id = $2", args.UserId, args.ProductId,
	)
	row.Scan(&result.Id, &result.Quantity)

	if result.Id == "" {
		fmt.Println("Product not found to update")
		utils.SendMsg("Product not found to update", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	_, err := trx.Exec("UPDATE cart SET quantity = $1 WHERE id = $2", result.Quantity+args.Quantity, result.Id)

	if err != nil {
		fmt.Println("Failed to update cart.")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendMsg("Cart updated", http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func GetCart(userId string, w http.ResponseWriter) {
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
			left join cart c on c.product_id  = p.id
			where c.customer_id = $1`,
		userId,
	)

	if err != nil {
		fmt.Println("Failed to get cart")
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

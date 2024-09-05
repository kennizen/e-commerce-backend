package service

import (
	"fmt"
	"net/http"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/utils"
)

type AddToCartArgs struct {
	UserId    string
	ProductId string
	Quantity  int
}

func AddToCart(args AddToCartArgs, w http.ResponseWriter) {
	var result struct {
		Id       string
		Quantity int
	}

	row := db.DB.QueryRow(
		"SELECT id, quantity FROM cart WHERE customer_id = $1 AND product_id = $2", args.UserId, args.ProductId,
	)
	row.Scan(&result.Id, &result.Quantity)

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	if result.Id == "" {
		_, err := trx.Exec(
			"INSERT INTO cart (customer_id, product_id, quantity) VALUES ($1, $2, $3)",
			args.UserId, args.ProductId, args.Quantity,
		)

		if err != nil {
			fmt.Println("Failed to add product to cart.")
			utils.SendMsg("Server error", http.StatusInternalServerError, w)
			return
		}

	} else {
		_, err := trx.Exec("UPDATE cart SET quantity = $1 WHERE id = $2", result.Quantity+args.Quantity, result.Id)

		if err != nil {
			fmt.Println("Failed to add product to cart.")
			utils.SendMsg("Server error", http.StatusInternalServerError, w)
			return
		}
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendMsg("Product added to cart", http.StatusOK, w)
}

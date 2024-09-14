package service

import (
	"fmt"
	"net/http"
	// "github.com/kennizen/e-commerce-backend/db"
)

type Product struct {
	ProductId string `validate:"required"`
	Quantity  int    `validate:"required,gte=1"`
}

type OrdersPayload struct {
	Products      []Product
	AddressUsed   string `validate:"required"`
	PaymentMethod string `validate:"required,oneof=ONLINE COD"`
}

func getIds(products *[]Product) string {
	// str := ""

	for _, prod := range *products {
		fmt.Println("id", prod.ProductId)
	}

	return ""
}

/*
WITH user_input AS (
    SELECT * FROM (values (1, 2), (2, 3), (3, 1)) AS t(product_id, quantity)
)
SELECT SUM(p.price * ui.quantity) AS total_order_price
FROM user_input ui
left join products p ON ui.product_id = p.id;
*/

func PlaceOrder(args OrdersPayload, userId string, w http.ResponseWriter) {
	getIds(&args.Products)
	// totalAmount := db.DB.QueryRow("SELECT SUM(p.price) FROM products p WHERE p.id IN " + getIds(&args.Products))
}

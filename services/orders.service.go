package service

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kennizen/e-commerce-backend/db"
	"github.com/kennizen/e-commerce-backend/utils"
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
	str := ""

	for _, prod := range *products {
		str = str + "(" + prod.ProductId + "," + " " + strconv.Itoa(prod.Quantity) + "),"
	}

	return strings.TrimSuffix(str, ",")
}

func getOrderValues(products *[]Product, userId string, paymentId string, addressId string) string {
	str := ""
	parentStr := ""

	for _, prod := range *products {
		str = str + "(" + userId + ", " + prod.ProductId + ", " + strconv.Itoa(prod.Quantity) + ", " + paymentId + ", " + addressId + ")"
		parentStr = parentStr + str + ","
		str = ""
	}

	return strings.TrimSuffix(parentStr, ",")
}

func PlaceOrder(args OrdersPayload, userId string, w http.ResponseWriter) {
	var amount float32

	row := db.DB.QueryRow(`
	WITH user_input AS (SELECT * FROM (values ` + getIds(&args.Products) + `)` + ` AS t(product_id, quantity)) 
	SELECT SUM(p.price * ui.quantity) FROM user_input ui LEFT JOIN products p ON ui.product_id = p.id`,
	)

	row.Scan(&amount)

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	row1 := trx.QueryRow(
		"INSERT INTO payments (payment_by, payment_method, payment_status, amount) VALUES ($1, $2, $3, $4) RETURNING id",
		userId, args.PaymentMethod, "success", amount,
	)

	var paymentId string
	row1.Scan(&paymentId)

	_, err1 := trx.Exec(`
	INSERT INTO orders (order_by, product_id, quantity, payment_id, address_used) 
	VALUES ` + getOrderValues(&args.Products, userId, paymentId, args.AddressUsed),
	)

	if err1 != nil {
		trx.Rollback()
		fmt.Println("Error inserting data in orders", err1.Error())
		utils.SendMsg("Bad request", http.StatusBadRequest, w)
		return
	}

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendMsg("Order placed", http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

type OrderRes struct {
	OrderId       string
	Quantity      float32
	OrderTime     string
	PaymentMethod string
	PaymentStatus string
	PaymentTime   string
}

type ProductRes struct {
	ProductId   string
	Title       string
	Description string
	Category    string
	Price       float32
	Image       string
	Thumbnail   string
}

type AddressRes struct {
	UserAddressPayload
}

type AllOrdersResponse struct {
	Order   OrderRes
	Product ProductRes
	Address AddressRes
}

func GetOrders(userId string, w http.ResponseWriter) {
	rows, err := db.DB.Query(`
		select
		o.id as order_id,
		o.quantity,
		o.created_at as order_time,
		p.id as product_id,
		p.title,
		p.description,
		p.category,
		p.price,
		p.image,
		p.thumbnail,
		p2.payment_method,
		p2.payment_status,
		p2.created_at as payment_time,
		a.country,
		a.state,
		a.address,
		a.zipcode,
		a.phone_number
		from orders o 
		left join products p on p.id = o.product_id 
		left join payments p2 on p2.id = o.payment_id 
		left join addresses a on a.id = o.address_used 
		where o.order_by = $1
		order by o.created_at
	`, userId)

	if err != nil {
		fmt.Println("Error in querying orders", err.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	defer rows.Close()

	var data AllOrdersResponse
	var resp []AllOrdersResponse
	isEmpty := true

	for rows.Next() {
		isEmpty = false
		err := rows.Scan(
			&data.Order.OrderId,
			&data.Order.Quantity,
			&data.Order.OrderTime,
			&data.Product.ProductId,
			&data.Product.Title,
			&data.Product.Description,
			&data.Product.Category,
			&data.Product.Price,
			&data.Product.Image,
			&data.Product.Thumbnail,
			&data.Order.PaymentMethod,
			&data.Order.PaymentStatus,
			&data.Order.PaymentTime,
			&data.Address.Country,
			&data.Address.State,
			&data.Address.Address,
			&data.Address.Zipcode,
			&data.Address.PhoneNo,
		)

		if err != nil {
			log.Fatalln("Error scanning row", err.Error())
		}

		resp = append(resp, data)
	}

	if isEmpty {
		utils.SendJson(utils.ResUserWithData{Msg: "Orders data", Data: make([]any, 0)}, http.StatusOK, w)
		return
	}

	utils.SendJson(utils.ResUserWithData{Msg: "Orders data", Data: resp}, http.StatusOK, w)
}

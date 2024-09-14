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

type UserDetailsPayload struct {
	Firstname  string `validate:"required"`
	Middlename string
	Lastname   string `validate:"required"`
	Age        int    `validate:"required,gte=1"`
	Email      string `validate:"required,email"`
	Avatar     string
}

type UserAddressPayload struct {
	Country string `validate:"required"`
	State   string `validate:"required"`
	Zipcode string `validate:"required"`
	PhoneNo string `validate:"required"`
	Address string `validate:"required"`
}

func UpdateUserDetails(args UserDetailsPayload, userId string, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM customers WHERE id = $1", userId)
	row.Scan(&id)

	if id == "" {
		fmt.Println("User not found to update")
		utils.SendMsg("User not found to update", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	trxRow := trx.QueryRow(
		"UPDATE customers SET firstname = $1, middlename = $2, lastname = $3, email = $4, age = $5, avatar = $6, updated_at = $7 WHERE id = $8 RETURNING *", args.Firstname, args.Middlename, args.Lastname, args.Email, args.Age, args.Avatar, time.Now().UTC().Format(time.RFC3339), id,
	)

	var user models.User

	trxRow.Scan(
		&user.Id,
		&user.Firstname,
		&user.Middlename,
		&user.Lastname,
		&user.Email,
		&user.Age,
		&user.Avatar,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(map[string]interface{}{
		"message": "User updated",
		"data":    user,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func DeleteUser(userId string, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM customers WHERE id = $1", userId)
	row.Scan(&id)

	if id == "" {
		fmt.Println("User not found to delete")
		utils.SendMsg("User not found to delete", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	trxRow := trx.QueryRow("DELETE FROM customers WHERE id = $1 RETURNING *", userId)

	var user models.User

	trxRow.Scan(
		&user.Id,
		&user.Firstname,
		&user.Middlename,
		&user.Lastname,
		&user.Email,
		&user.Age,
		&user.Avatar,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(map[string]interface{}{
		"message": "User deleted",
		"data":    user,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func AddAddress(args UserAddressPayload, userId string, w http.ResponseWriter) {
	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	row := trx.QueryRow(
		"INSERT INTO addresses (country, state, address, zipcode, phone_number, address_of) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *",
		args.Country, args.State, args.Address, args.Zipcode, args.PhoneNo, userId,
	)

	var address models.Address

	row.Scan(
		&address.Id,
		&address.Country,
		&address.State,
		&address.Address,
		&address.Zipcode,
		&address.PhoneNumber,
		&address.AddressOf,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(map[string]interface{}{
		"message": "Address added",
		"data":    address,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func UpdateAddress(args UserAddressPayload, addressId string, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM addresses WHERE id = $1", addressId)
	row.Scan(&id)

	if id == "" {
		fmt.Println("Address not found to update")
		utils.SendMsg("Address not found to update", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	addr := trx.QueryRow(
		`UPDATE addresses 
		SET country = $1, state = $2, address = $3, zipcode = $4, phone_number = $5, updated_at = $6 
		WHERE id = $7 RETURNING *`,
		args.Country, args.State, args.Address, args.Zipcode, args.PhoneNo, time.Now().UTC().Format(time.RFC3339), addressId,
	)

	var address models.Address

	addr.Scan(
		&address.Id,
		&address.Country,
		&address.State,
		&address.Address,
		&address.Zipcode,
		&address.PhoneNumber,
		&address.AddressOf,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(map[string]interface{}{
		"message": "Address updated",
		"data":    address,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func DeleteAddress(addressId string, userId string, w http.ResponseWriter) {
	var id string = ""

	row := db.DB.QueryRow("SELECT id FROM addresses WHERE id = $1", addressId)
	row.Scan(&id)

	if id == "" {
		fmt.Println("Address not found to delete")
		utils.SendMsg("Address not found to delete", http.StatusBadRequest, w)
		return
	}

	trx, trxErr := db.DB.Begin()

	if trxErr != nil {
		fmt.Println("Error in transaction", trxErr.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	addr := trx.QueryRow(`DELETE FROM addresses WHERE id = $1 AND address_of = $2 RETURNING *`, addressId, userId)

	var address models.Address

	addr.Scan(
		&address.Id,
		&address.Country,
		&address.State,
		&address.Address,
		&address.Zipcode,
		&address.PhoneNumber,
		&address.AddressOf,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	comErr := trx.Commit()

	if comErr != nil {
		fmt.Println("Error in transaction")
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	utils.SendJson(map[string]interface{}{
		"message": "Address deleted",
		"data":    address,
	}, http.StatusOK, w)
}

// ---------------------------------------------------------------------------------------- //

func GetAddresses(userId string, w http.ResponseWriter) {
	rows, err := db.DB.Query("SELECT * FROM addresses WHERE address_of = $1", userId)

	if err != nil {
		fmt.Println("Failed to query addresses", err.Error())
		utils.SendMsg("Server error", http.StatusInternalServerError, w)
		return
	}

	defer rows.Close()

	var address models.Address
	var addresses []models.Address
	isEmpty := true

	for rows.Next() {
		isEmpty = false
		err := rows.Scan(
			&address.Id,
			&address.Country,
			&address.State,
			&address.Address,
			&address.Zipcode,
			&address.PhoneNumber,
			&address.AddressOf,
			&address.CreatedAt,
			&address.UpdatedAt,
		)

		if err != nil {
			log.Fatalln("Error scanning row", err.Error())
		}

		addresses = append(addresses, address)
	}

	if isEmpty {
		utils.SendJson(map[string][]any{"data": make([]any, 0)}, http.StatusOK, w)
		return
	}

	utils.SendJson(map[string]any{"data": addresses}, http.StatusOK, w)
}

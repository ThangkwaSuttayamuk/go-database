package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func getUsers() ([]User, error) {
	rows, err := db.Query("SELECT id, firstname, lastname FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Firstname, &u.Lastname)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func createProduct(product *Product) error {
	currentTime := (time.Now()).Format("2006-01-02")

	_, err := db.Exec(
		"INSERT INTO public.products(name, price, supplier_id, time) VALUES ($1,$2,$3,$4);",
		product.Name, product.Price, product.Supplier_id, currentTime,
	)

	return err
}

func createUser(user *User) error {
	_, err := db.Exec(
		"INSERT INTO public.users(firstname, lastname) VALUES ($1, $2);",
		user.Firstname, user.Lastname,
	)

	return err
}

func addDayOff(uid int, dayOff *DayOff) error {
	_, err := db.Exec(
		"INSERT INTO public.day_off(date, uid) VALUES ($1,$2);",
		dayOff.Date, uid,
	)

	return err
}

func deleteDayOff(uid int, dayOff *DayOff) error {
	result, err := db.Exec(
		"DELETE FROM public.day_off WHERE uid = $1 AND date = $2;",
		uid, dayOff.Date,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no record found to delete")
	}

	return nil
}

func getProductById(id int) (Product, error) {
	var p Product

	row := db.QueryRow(
		"SELECT id, name, price FROM public.products WHERE id = $1;",
		id,
	)

	err := row.Scan(&p.ID, &p.Name, &p.Price)

	// Handle the case where no rows were found
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a custom error indicating no product found
			return Product{}, fmt.Errorf("no product found with id %d", id)
		}
		// Return any other errors encountered during scanning
		return Product{}, err
	}

	return p, err
}

//choice 1
// func updateProduct(id int, product *Product) error {
// 	_, err := db.Exec(
// 		"UPDATE public.products SET name = $1, price = $2 WHERE id = $3;",
// 		product.Name, product.Price, id,
// 	)

// 	return err
// }

// choice 2
func updateProduct(id int, product *Product) (Product, error) {

	var p Product

	row := db.QueryRow(
		"UPDATE public.products SET name = $1, price = $2 WHERE id = $3 RETURNING id, name, price;",
		product.Name, product.Price, id,
	)

	err := row.Scan(&p.ID, &p.Name, &p.Price)

	if err != nil {
		return Product{}, err
	}

	return p, err
}

func deleteProduct(id int) error {
	_, err := db.Exec(
		"DELETE FROM public.products WHERE id = $1;",
		id,
	)

	return err
}

func getProducts() ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func getDayOffs(uid int) ([]DayOff, error) {
	rows, err := db.Query("SELECT id, date, uid FROM day_off WHERE uid = $1;", uid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var date []DayOff

	for rows.Next() {
		var d DayOff
		err := rows.Scan(&d.ID, &d.Date, &d.Uid)
		if err != nil {
			return nil, err
		}
		date = append(date, d)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return date, nil
}

func addProductAndSupplier(supplierName string, productName string, productPrice int) error {

	var s Supplier

	// Start a transaction
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	// Rollback the transaction in case of a panic
	defer tx.Rollback()

	// Insert into the supplier table
	supplierResult := tx.QueryRow("INSERT INTO suppliers (name) VALUES ($1) RETURNING id, name;", supplierName)

	err = supplierResult.Scan(&s.ID, &s.Name)

	if err != nil {
		return err
	}

	// Insert into the product table
	_, err = tx.Exec("INSERT INTO products (name, price, supplier_id) VALUES ($1, $2, $3);",
		productName, productPrice, s.ID,
	)

	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

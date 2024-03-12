package repository

import (
	"database/sql"
	"example.com/server/pkg/models"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db         *sql.DB
	dbInitOnce sync.Once
)

func ConnectDB() (*sql.DB, error) {
	connectionInfo := "host=localhost port=5432 user=postgres password=12345 dbname=mydatabase sslmode=disable"

	db, err := sql.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func getDB() (*sql.DB, error) {
	var err error

	dbInitOnce.Do(func() {
		db, err = ConnectDB()
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertUser(userData models.User) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", userData.Username, userData.Email, userData.Password)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User

	db, err := getDB()
	if err != nil {
		return user, err
	}

	err = db.QueryRow("SELECT id, username, email, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func InsertProduct(userID int, productData models.Product) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO products (user_id, name, price) VALUES ($1, $2, $3)", userID, productData.Name, productData.Price)
	if err != nil {
		return err
	}

	return nil
}

func GetProducts() ([]models.Product, error) {
	var products []models.Product

	db, err := getDB()
	if err != nil {
		return products, err
	}

	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return products, err
		}
		products = append(products, product)
	}

	return products, nil
}

func GetProductByID(id int) (models.Product, error) {
	var product models.Product

	db, err := getDB()
	if err != nil {
		return product, err
	}

	err = db.QueryRow("SELECT id, name, price FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return product, err
	}

	return product, nil
}

func UpdateProduct(id int, newData models.Product) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", newData.Name, newData.Price, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteProduct(id int) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

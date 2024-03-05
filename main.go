package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type RegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Структура для успешного ответа
type SuccessResponse struct {
	Message string `json:"message"`
}

// Структура для ответа с данными пользователя
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Структура для ошибочного ответа
type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	// Инициализация маршрутизатора Gin
	router := gin.Default()
	store := cookie.NewStore([]byte("your_secret_key_here"))
	router.Use(sessions.Sessions("mysession", store))

	// Подключение к базе данных PostgreSQL
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Маршрут для регистрации новых пользователей
	router.POST("/register", func(c *gin.Context) {
		var req RegistrationRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Хэширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to hash password"})
			return
		}

		user := User{Username: req.Username, Email: req.Email, Password: string(hashedPassword)}

		// Вставка данных пользователя в таблицу
		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to register user"})
			return
		}

		// Получаем сеанс
		session := sessions.Default(c)
		// Устанавливаем пользовательский идентификатор в сеанс
		session.Set("userID", user.ID)
		// Сохраняем сеанс
		session.Save()

		c.JSON(http.StatusOK, SuccessResponse{Message: "User registered successfully"})
	})

	router.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}
		var user User
		err := db.QueryRow("SELECT id, username, email, password FROM users WHERE username = $1", req.Username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
			return
		}

		// Получаем сеанс
		session := sessions.Default(c)
		// Устанавливаем пользовательский идентификатор в сеанс
		session.Set("userID", user.ID)
		// Сохраняем сеанс
		session.Save()

		c.JSON(http.StatusOK, SuccessResponse{Message: "Login successful"})
	})

	// Маршрут для получения информации о сессии пользователя
	router.GET("/session-info", func(c *gin.Context) {
		// Получаем сеанс
		session := sessions.Default(c)
		// Получаем значение из сессии (например, идентификатор пользователя)
		userID := session.Get("userID")

		// Проверяем, есть ли значение в сессии
		if userID == nil {
			c.JSON(http.StatusOK, ErrorResponse{Error: "Session not found"})
			return
		}

		// Если значение найдено, отображаем его
		c.JSON(http.StatusOK, SuccessResponse{Message: fmt.Sprintf("User ID from session: %v", userID)})
	})

	// Маршрут для создания нового товара
	router.POST("/products", func(c *gin.Context) {
		var product Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Вставка данных о товаре в таблицу
		_, err := db.Exec("INSERT INTO products (name, price) VALUES ($1, $2)", product.Name, product.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create product"})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{Message: "Product created successfully"})
	})

	// Маршрут для получения списка всех товаров
	router.GET("/products", func(c *gin.Context) {
		var products []Product
		rows, err := db.Query("SELECT id, name, price FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch products"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var product Product
			if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to scan product"})
				return
			}
			products = append(products, product)
		}

		c.JSON(http.StatusOK, products)
	})

	// Маршрут для получения информации о конкретном товаре по его ID
	router.GET("/products/:id", func(c *gin.Context) {
		var product Product
		id := c.Param("id")

		row := db.QueryRow("SELECT id, name, price FROM products WHERE id = $1", id)
		if err := row.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	})

	// Маршрут для обновления информации о товаре по его ID
	router.PUT("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var product Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}

		_, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", product.Name, product.Price, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update product"})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{Message: "Product updated successfully"})
	})

	// Маршрут для удаления товара по его ID
	router.DELETE("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM products WHERE id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete product"})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{Message: "Product deleted successfully"})
	})

	// Запуск сервера на порту 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Функция для создания подключения к базе данных PostgreSQL
func connectDB() (*sql.DB, error) {
	connectionInfo := "host=localhost port=5432 user=postgres password=12345 dbname=mydatabase sslmode=disable"

	db, err := sql.Open("postgres", connectionInfo)
	if err != nil {

		return nil, err
	}

	// Проверка подключения к базе данных
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

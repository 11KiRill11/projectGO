package handlers

import (
	"example.com/server/pkg/models"
	"example.com/server/pkg/services"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	var req models.RegistrationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	// Создание экземпляра структуры User из RegistrationRequest
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// Вызов нужной функции из services/auth.go
	err := services.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "User registered successfully"})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	_, err := services.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Login successful"})
}

func GetSessionInfo(c *gin.Context) {

	session := sessions.Default(c)

	userID := session.Get("userID")

	if userID == nil {
		c.JSON(http.StatusOK, models.ErrorResponse{Error: "Session not found"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: fmt.Sprintf("User ID from session: %v", userID)})
}

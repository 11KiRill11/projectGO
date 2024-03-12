package handlers

import (
	"example.com/server/pkg/models"
	"example.com/server/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	userID, _ := c.Get("userID")

	err := services.CreateProduct(userID.(int), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Product created successfully"})
}

func GetProducts(c *gin.Context) {
	products, err := services.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")

	// Преобразование строки в целое число
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	product, err := services.GetProductByID(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	// Преобразование строки в целое число
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	err = services.UpdateProduct(productID, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Product updated successfully"})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Преобразование строки в целое число
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	err = services.DeleteProduct(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Product deleted successfully"})
}

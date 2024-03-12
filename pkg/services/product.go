package services

import (
	"example.com/server/pkg/models"
	"example.com/server/pkg/repository"
)

func CreateProduct(userID int, productData models.Product) error {
	// Логика создания продукта
	err := repository.InsertProduct(userID, productData)
	if err != nil {
		return err
	}
	return nil
}

func GetProducts() ([]models.Product, error) {
	// Логика получения списка всех товаров
	products, err := repository.GetProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func GetProductByID(id int) (models.Product, error) {
	// Логика получения информации о конкретном товаре по его ID
	product, err := repository.GetProductByID(id)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func UpdateProduct(id int, newData models.Product) error {
	// Логика обновления информации о товаре по его ID
	err := repository.UpdateProduct(id, newData)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProduct(id int) error {
	// Логика удаления товара по его ID
	err := repository.DeleteProduct(id)
	if err != nil {
		return err
	}
	return nil
}

package services

import (
	"example.com/server/pkg/models"
	"example.com/server/pkg/repository"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(userData models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userData.Password = string(hashedPassword)
	err = repository.InsertUser(userData)
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(username, password string) (models.User, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return models.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

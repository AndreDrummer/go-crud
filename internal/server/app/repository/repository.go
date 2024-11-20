package repository

import (
	"encoding/json"
	"errors"
	dbhandler "go-crud/internal/server/app/db/handler"
	"go-crud/internal/server/app/model"
)

func GetUserList() ([]model.User, error) {
	db := dbhandler.OpenDB()
	data, err := db.FindAll()

	if err != nil {
		return []model.User{}, err
	}

	users := make([]model.User, len(data))

	for i, v := range data {
		if err := json.Unmarshal([]byte(v), &users[i]); err != nil {
			return users, errors.New("error converting the data from DB to JSON")
		}
	}

	return users, nil
}

func GetUser(userID string) (model.User, error) {

	db := dbhandler.OpenDB()
	userString, err := db.FindByID(userID)

	if err != nil {
		return model.User{}, err
	}

	var user model.User
	if err := json.Unmarshal([]byte(userString), &user); err != nil {
		return model.User{}, errors.New("error converting the data from DB to JSON")
	}

	return user, nil
}

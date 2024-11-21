package repository

import (
	"encoding/json"
	"fmt"
	dbhandler "go-crud/internal/server/app/db/handler"
	customerrors "go-crud/internal/server/app/errors"
	"go-crud/internal/server/app/model"

	"github.com/google/uuid"
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
			return users, &customerrors.JsonDecodingError{
				Type: fmt.Sprintf("%T", users),
				Err:  err,
			}
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
		return model.User{}, &customerrors.JsonDecodingError{
			Type: fmt.Sprintf("%T", user),
			Err:  err,
		}
	}

	return user, nil
}

func InsertUser(user *model.User) error {
	intUUID := uuid.New()
	userID := intUUID.String()
	user.ID = userID

	userJson, err := json.Marshal(user)
	if err != nil {
		return &customerrors.JsonEncodingError{
			Type: fmt.Sprintf("%T", user),
			Err:  err,
		}
	}

	db := dbhandler.OpenDB()

	if err := db.Insert(string(userJson)); err != nil {
		return err
	}

	return nil
}

func UpdateUser(user *model.User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return &customerrors.JsonEncodingError{
			Type: fmt.Sprintf("%T", user),
			Err:  err,
		}
	}

	db := dbhandler.OpenDB()

	if err := db.Update(user.ID, string(userJson)); err != nil {
		return err
	}

	return nil
}

func DeleteUser(userID string) error {
	db := dbhandler.OpenDB()

	if err := db.Delete(userID); err != nil {
		return err
	} else {
		return nil
	}
}

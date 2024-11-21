package repository

import (
	"encoding/json"
	"fmt"
	dbhandler "go-crud/internal/server/app/db/handler"
	customerrors "go-crud/internal/server/app/errors"
	"go-crud/internal/server/app/model"

	"github.com/google/uuid"
)

var volatile map[string]string

func insertUsersOnVolatileMemory(userID, userData string) {
	volatile[userID] = userData
}

func getAllUsersFromVolatileMemory() []string {
	var userData []string = []string{}

	for _, v := range volatile {
		userData = append(userData, v)
	}

	return userData
}

func getUserFromVolatileMemory(userID string) string {
	return volatile[userID]
}

func deleteUserFromVolatileMemory(userID string) bool {
	var removeSuccessfully bool

	if _, ok := volatile[userID]; ok {
		delete(volatile, userID)
		removeSuccessfully = true
	}

	return removeSuccessfully
}

func updateUserOnVolatileMemory(userID, userData string) {
	if _, ok := volatile[userID]; ok {
		volatile[userID] = userData
	}
}

func GetUserList() ([]model.User, error) {
	var data []string

	memData := getAllUsersFromVolatileMemory()

	if len(memData) == 0 {
		db := dbhandler.OpenDB()
		dbData, err := db.FindAll()

		if err != nil {
			return []model.User{}, err
		}

		data = dbData
	} else {
		data = memData
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
	var userData string

	memData := getUserFromVolatileMemory(userID)

	if len(memData) == 0 {
		db := dbhandler.OpenDB()
		dbData, err := db.FindByID(userID)

		if err != nil {
			return model.User{}, err
		}

		userData = dbData
	} else {
		userData = memData
	}

	var user model.User
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
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

	insertUsersOnVolatileMemory(user.ID, string(userJson))
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

	updateUserOnVolatileMemory(user.ID, string(userJson))

	return nil
}

func DeleteUser(userID string) error {

	deleteUserFromVolatileMemory(userID)

	db := dbhandler.OpenDB()

	if err := db.Delete(userID); err != nil {
		return err
	} else {
		return nil
	}
}

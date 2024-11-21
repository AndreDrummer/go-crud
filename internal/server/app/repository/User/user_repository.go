package userrepository

import (
	"encoding/json"
	"fmt"
	dbhandler "go-crud/internal/server/app/db/handler"
	customerrors "go-crud/internal/server/app/errors"
	"go-crud/internal/server/app/model"
	"go-crud/internal/server/app/repository"
	"log/slog"
	"regexp"

	"github.com/google/uuid"
)

type UserRepository struct {
	isInitialized bool
	volatileMem   map[string]string
}

func getEntryID(entry string) string {
	if entry == "" {
		return ""
	}

	re := regexp.MustCompile(`"([a-f0-9-]{36})"`)
	match := re.FindStringSubmatch(entry)

	if len(match) > 1 {
		return match[1]
	} else {
		return ""
	}
}

func getAllUsersFromPersistentMemory() ([]string, error) {
	db := dbhandler.OpenDB()
	dbData, err := db.FindAll()

	if err != nil {
		return []string{}, err
	}

	return dbData, nil
}

func (u *UserRepository) loadDataFromPersistentMemory() {
	dataFromDB, err := getAllUsersFromPersistentMemory()

	if err != nil {
		slog.Error("error when loading data from DB")
		return
	}

	for _, v := range dataFromDB {
		vID := getEntryID(v)
		if vID != "" {
			u.insertUsersOnVolatileMemory(vID, v)
		}
	}
}

func NewUserRepository() *UserRepository {
	ur := &UserRepository{
		isInitialized: true,
		volatileMem:   make(map[string]string),
	}

	ur.loadDataFromPersistentMemory()

	return ur
}

func (u *UserRepository) insertUsersOnVolatileMemory(userID, userData string) {
	u.volatileMem[userID] = userData
}

func (u *UserRepository) getAllUsersFromVolatileMemory() []string {
	var userData []string = []string{}

	for _, v := range u.volatileMem {
		userData = append(userData, v)
	}

	return userData
}

func (u *UserRepository) getUserFromVolatileMemory(userID string) string {
	return u.volatileMem[userID]
}

func (u *UserRepository) deleteUserFromVolatileMemory(userID string) bool {
	var removeSuccessfully bool

	if _, ok := u.volatileMem[userID]; ok {
		delete(u.volatileMem, userID)
		removeSuccessfully = true
	}

	return removeSuccessfully
}

func (u *UserRepository) updateUserOnVolatileMemory(userID, userData string) {
	if _, ok := u.volatileMem[userID]; ok {
		u.volatileMem[userID] = userData
	}
}

func (u *UserRepository) GetAll() ([]repository.Document, error) {
	var data []string

	memData := u.getAllUsersFromVolatileMemory()

	if len(memData) == 0 {
		dbData, err := getAllUsersFromPersistentMemory()

		if err != nil {
			return []repository.Document{}, err
		}

		data = dbData
	} else {
		data = memData
	}

	users := make([]model.User, len(data))

	for i, v := range data {
		if err := json.Unmarshal([]byte(v), &users[i]); err != nil {
			return []repository.Document{}, &customerrors.JsonDecodingError{
				Type: fmt.Sprintf("%T", users),
				Err:  err,
			}
		}
	}

	entities := make([]repository.Document, len(users))

	for i, v := range users {
		entities[i] = repository.Document(v)
	}

	return entities, nil
}

func (u *UserRepository) GetOne(userID string) (repository.Document, error) {
	var userData string

	memData := u.getUserFromVolatileMemory(userID)

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

func (u *UserRepository) Insert(document repository.Document) error {
	user, ok := document.(*model.User)

	if !ok {
		return fmt.Errorf("invalid document type: expected *model.User, received %T", document)
	}

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

	u.insertUsersOnVolatileMemory(user.ID, string(userJson))
	return nil
}
func (u *UserRepository) Update(document repository.Document) error {
	user, ok := document.(*model.User)

	if !ok {
		return fmt.Errorf("invalid document type: expected *model.User, received %T", document)
	}

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

	u.updateUserOnVolatileMemory(user.ID, string(userJson))

	return nil
}

func (u *UserRepository) Delete(userID string) error {

	u.deleteUserFromVolatileMemory(userID)

	db := dbhandler.OpenDB()

	if err := db.Delete(userID); err != nil {
		return err
	} else {
		return nil
	}
}

package model

import (
	"fmt"
	"strings"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

func (u *User) IsValid() bool {
	return u.FirstName != "" && strings.TrimSpace(u.FirstName) != "" &&
		u.Biography != "" && strings.TrimSpace(u.Biography) != "" &&
		u.LastName != "" && strings.TrimSpace(u.LastName) != ""
}

func (u User) ToString() string {
	return fmt.Sprintf("Document of type %T", u)
}

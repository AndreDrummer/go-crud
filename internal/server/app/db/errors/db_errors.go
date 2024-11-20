package dberrors

import "fmt"

type DBNotFoundError struct{}

func (d *DBNotFoundError) Error() string {
	return "Entity Not Found"
}

type DBPathError struct {
	Err error
}

func (d *DBPathError) Error() string {
	return fmt.Sprintf("PathError: %s", d.Err.Error())
}

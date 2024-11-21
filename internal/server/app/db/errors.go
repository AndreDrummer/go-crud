package db

import "fmt"

type DBPathError struct {
	Err error
}

func (d *DBPathError) Error() string {
	return fmt.Sprintf("PathError: %s", d.Err.Error())
}

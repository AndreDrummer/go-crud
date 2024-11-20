package dberrors

type DBNotFoundError struct{}

func (d *DBNotFoundError) Error() string {
	return "ID Not Found"
}

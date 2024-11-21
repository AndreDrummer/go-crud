package customerrors

import "fmt"

type NotFoundError struct {
	Msg string
}

func (n *NotFoundError) Error() string {
	if n.Msg == "" {
		return "Entity Not Found"
	} else {
		return n.Msg
	}
}

type JsonEncodingError struct {
	Type any
	Err  error
}

func (j *JsonEncodingError) Error() string {
	return fmt.Sprintf("Error converting data of type %v to JSON, %s", j.Type, j.Err.Error())

}

type JsonDecodingError struct {
	Type any
	Err  error
}

func (j *JsonDecodingError) Error() string {
	return fmt.Sprintf("Error converting JSON to %v: %s", j.Type, j.Err.Error())
}

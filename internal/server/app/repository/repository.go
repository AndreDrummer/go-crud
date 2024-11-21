package repository

type Document interface {
	ToString() string
}

type Repository interface {
	CheckInitialization() error
	GetAll() ([]Document, error)
	GetOne(documentID string) (Document, error)
	Insert(Document) error
	Update(Document) error
	Delete(documentID string) error
}

const (
	User = "user"
)

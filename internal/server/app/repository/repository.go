package repository

type Entity interface {
	ToString() string
}

type Repository interface {
	GetAll() ([]Entity, error)
	GetOne(entityID string) (Entity, error)
	Insert(Entity) error
	Update(Entity) error
	Delete(entityID string) error
}

const (
	User = "user"
)

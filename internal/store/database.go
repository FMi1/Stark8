package store

type Database interface {
	GetUser(string) (User, error)
	CreateUser(UserParams) error
}

package users

type UserRepository interface {
	CreateUser(*User) (*User, error)
	UpdateUser(*User) (*User, error)
	ReadUser(int) (*User, error)
	ListUsers() ([]User, error)
}

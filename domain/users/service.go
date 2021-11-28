package users

// UserService defines author service behavior.
type UserService interface {
	CreateUser(*User) (*User, error)
	ReadUser(id int) (*User, error)
	UpdateUser(*User) (*User, error)
	ListUsers() ([]User, error)
}

// Service struct handles author business logic tasks.
type Service struct {
	repository UserRepository
}

func (svc *Service) CreateUser(user *User) (*User, error) {
	return svc.repository.CreateUser(user)
}

func (svc *Service) ReadUser(id int) (*User, error) {
	return svc.repository.ReadUser(id)
}
func (svc *Service) UpdateUser(user *User) (*User, error) {
	return svc.repository.UpdateUser(user)
}

func (svc *Service) ListUsers() ([]User, error) {
	return svc.repository.ListUsers()
}

// NewService creates a new service struct
func NewService(repository UserRepository) *Service {
	return &Service{repository: repository}
}

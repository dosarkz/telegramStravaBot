package users

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	domainErrors "telegramStravaBot/app/handlers"
)

const (
	createError = "Error in creating new user"
	readError   = "Error in finding user in the database"
	listError   = "Error in getting workout from the database"
)

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserRepository {
	db.AutoMigrate(&User{})
	return &UserRepository{db}
}

func (s *UserRepository) CreateUser(user *User) (*User, error) {
	entity := toDBModel(user)

	if err := s.db.Create(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, createError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toDomainModel(entity), nil
}

func (s *UserRepository) UpdateUser(user *User) (*User, error) {
	result := &User{}
	query := s.db.Model(&result).Updates(user)

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, readError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toDomainModel(result), nil
}

func (s *UserRepository) ReadUser(id int) (*User, error) {
	result := &User{}

	query := s.db.Where("id = ?", id).First(result)

	if query.RecordNotFound() {
		appErr := domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		return nil, appErr
	}

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, readError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toDomainModel(result), nil
}

func (s *UserRepository) ListUsers() ([]User, error) {
	var results []User

	if err := s.db.Find(&results).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, listError), domainErrors.RepositoryError)
		return nil, appErr
	}

	var users = make([]User, len(results))

	for i, element := range results {
		users[i] = *toDomainModel(&element)
	}

	return users, nil
}

func (s *UserRepository) FindUserByTelegramId(id int64) (*User, error) {
	result := &User{}

	query := s.db.Where("telegram_id = ?", id).First(result)

	if query.RecordNotFound() {
		appErr := domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		return nil, appErr
	}

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, readError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toDomainModel(result), nil
}

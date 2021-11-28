package users

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	domainErrors "telegramStravaBot/domain/errors"
	domain "telegramStravaBot/domain/users"
)

const (
	createError = "Error in creating new user"
	readError   = "Error in finding user in the database"
	listError   = "Error in getting candidate from the database"
)

// Store struct manages interactions with candidate store
type Store struct {
	db *gorm.DB
}

// New creates a new Store struct
func New(db *gorm.DB) *Store {
	db.AutoMigrate(&User{})

	return &Store{
		db: db,
	}
}

func (s *Store) CreateUser(user *domain.User) (*domain.User, error) {
	entity := toDBModel(user)

	if err := s.db.Create(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, createError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toDomainModel(entity), nil
}

func (s *Store) UpdateUser(user *domain.User) (*domain.User, error) {
	result := &User{}
	query := s.db.Model(&result).Updates(user)

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, readError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toDomainModel(result), nil
}

func (s *Store) ReadUser(id int) (*domain.User, error) {
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

func (s *Store) ListUsers() ([]domain.User, error) {
	var results []User

	if err := s.db.Find(&results).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, listError), domainErrors.RepositoryError)
		return nil, appErr
	}

	var users = make([]domain.User, len(results))

	for i, element := range results {
		users[i] = *toDomainModel(&element)
	}

	return users, nil
}

func (s *Store) FindByAppId(id int) (*domain.User, error) {
	result := &User{}

	query := s.db.Where("app_id = ?", id).First(result)

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

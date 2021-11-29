package workouts

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	domainErrors "telegramStravaBot/domain/errors"
	domain "telegramStravaBot/domain/workouts"
)

const (
	createError = "Error in creating new user"
	readError   = "Error in finding user in the database"
	listError   = "Error in getting workout from the database"
)

// Store struct manages interactions with workout store
type Store struct {
	db *gorm.DB
}

// New creates a new Store struct
func New(db *gorm.DB) *Store {
	db.AutoMigrate(&Workout{})

	return &Store{
		db: db,
	}
}

func (s *Store) CreateWorkout(workout *domain.Workout) (*domain.Workout, error) {
	entity := toDBModel(workout)

	if err := s.db.Create(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, createError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toDomainModel(entity), nil
}

func (s *Store) UpdateWorkout(workout *domain.Workout) (*domain.Workout, error) {
	result := &Workout{}
	query := s.db.Model(&result).Updates(workout)

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, readError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toDomainModel(result), nil
}

func (s *Store) ReadWorkout(id int) (*domain.Workout, error) {
	result := &Workout{}

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

func (s *Store) ListWorkouts() ([]domain.Workout, error) {
	var results []Workout

	if err := s.db.Find(&results).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, listError), domainErrors.RepositoryError)
		return nil, appErr
	}

	var workouts = make([]domain.Workout, len(results))

	for i, element := range results {
		workouts[i] = *toDomainModel(&element)
	}

	return workouts, nil
}

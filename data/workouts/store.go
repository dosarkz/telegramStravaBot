package workouts

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	domainErrors "telegramStravaBot/domain/errors"
	domain "telegramStravaBot/domain/workouts"
	"time"
)

const (
	createError = "Error in creating new user"
	readError   = "Error in finding user in the database"
	listError   = "Error in getting workout from the database"
)

// Store struct manages interactions with workout store
type Store struct {
	db            *gorm.DB
	workoutUserDb *gorm.DB
}

// New creates a new Store struct
func New(db *gorm.DB, workoutUserDb *gorm.DB) *Store {
	db.AutoMigrate(&Workout{})
	workoutUserDb.AutoMigrate(&WorkoutUser{})

	return &Store{
		db:            db,
		workoutUserDb: workoutUserDb,
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
	var currentTime = time.Now()

	if err := s.db.Where("created_at >= ?", currentTime).Find(&results).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, listError), domainErrors.RepositoryError)
		return nil, appErr
	}

	var workouts = make([]domain.Workout, len(results))

	for i, element := range results {
		workouts[i] = *toDomainModel(&element)
	}

	return workouts, nil
}

func (s *Store) ListWorkoutMembers(workoutId int) ([]domain.WorkoutUser, error) {
	var results []WorkoutUser

	if err := s.db.Where("workout_id = ?", workoutId).Find(&results).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, listError), domainErrors.RepositoryError)
		return nil, appErr
	}

	var workouts = make([]domain.WorkoutUser, len(results))

	for i, element := range results {
		workouts[i] = *toWorkoutUserDomainModel(&element)
	}

	return workouts, nil
}

func (s *Store) FindBy(userId int, workoutId int) (*domain.WorkoutUser, error) {
	result := &WorkoutUser{}
	query := s.workoutUserDb.Where("user_id = ? and workout_id = ?", userId, workoutId).First(result)

	if query.RecordNotFound() {
		appErr := domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		return nil, appErr
	}

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, readError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toWorkoutUserDomainModel(result), nil
}

func (s *Store) CreateWorkoutUser(workoutUser *domain.WorkoutUser) (*domain.WorkoutUser, error) {
	entity := toWorkoutUserDBModel(workoutUser)

	if err := s.db.Create(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, createError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toWorkoutUserDomainModel(entity), nil
}

func (s *Store) Delete(workoutUser *domain.WorkoutUser) (bool, error) {

	entity := toWorkoutUserDBModel(workoutUser)
	query := s.workoutUserDb.Delete(entity)

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, readError), domainErrors.RepositoryError)
		return false, appErr
	}

	return true, nil
}

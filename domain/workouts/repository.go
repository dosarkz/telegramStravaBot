package workouts

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	domainErrors "telegramStravaBot/app/handlers"
	"telegramStravaBot/domain/users"
	"time"
)

const (
	createError = "Error in creating new user"
	readError   = "Error in finding user in the database"
	listError   = "Error in getting workout from the database"
)

type WorkoutRepository struct {
	db *gorm.DB
}

// New creates a new Store struct
func New(db *gorm.DB) *WorkoutRepository {
	db.AutoMigrate(&Workout{})

	return &WorkoutRepository{
		db: db,
	}
}

func (s *WorkoutRepository) CreateWorkout(workout *Workout) (*Workout, error) {
	entity := toDBModel(workout)

	if err := s.db.Create(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, createError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toDomainModel(entity), nil
}

func (s *WorkoutRepository) UpdateWorkout(workout *Workout) (*Workout, error) {
	result := &Workout{}
	query := s.db.Model(&result).Updates(workout)

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, readError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toDomainModel(result), nil
}

func (s *WorkoutRepository) ReadWorkout(id int) (*Workout, error) {
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

func (s *WorkoutRepository) ListWorkouts() ([]Workout, error) {
	var results []Workout
	var currentTime = time.Now()

	if err := s.db.Where("created_at >= ?", currentTime).Find(&results).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, listError), domainErrors.RepositoryError)
		return nil, appErr
	}

	var workouts = make([]Workout, len(results))

	for i, element := range results {
		workouts[i] = *toDomainModel(&element)
	}

	return workouts, nil
}

func (s *WorkoutRepository) ListWorkoutMembers(workoutId int) ([]WorkoutUser, error) {
	var results []WorkoutUser

	if err := s.db.Where("workout_id = ?", workoutId).Find(&results).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, listError), domainErrors.RepositoryError)
		return nil, appErr
	}

	var workouts = make([]WorkoutUser, len(results))

	for i, element := range results {
		workouts[i] = *toWorkoutUserDomainModel(&element)
	}

	return workouts, nil
}

func (s *WorkoutRepository) FindBy(userId int, workoutId int) (*WorkoutUser, error) {
	result := &WorkoutUser{}
	query := s.db.Where("user_id = ? and workout_id = ?", userId, workoutId).First(result)

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

func (s *WorkoutRepository) CreateWorkoutUser(workoutUser *WorkoutUser) (*WorkoutUser, error) {
	entity := toWorkoutUserDBModel(workoutUser)

	if err := s.db.Create(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, createError), domainErrors.RepositoryError)
		return nil, appErr
	}

	return toWorkoutUserDomainModel(entity), nil
}

func (s *WorkoutRepository) Delete(workoutUser *WorkoutUser) (bool, error) {

	entity := toWorkoutUserDBModel(workoutUser)
	query := s.db.Delete(entity)

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, readError), domainErrors.RepositoryError)
		return false, appErr
	}

	return true, nil
}

func (s *WorkoutRepository) DeleteWorkout(workout *Workout) (bool, error) {
	entity := toDBModel(workout)
	query := s.db.Delete(entity)

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(errors.Wrap(err, readError), domainErrors.RepositoryError)
		return false, appErr
	}

	return true, nil
}

func (s WorkoutRepository) CallbackNewWorkout(update tgbotapi.Update) (bool, *Workout) {

	newText := strings.Split(update.Message.Text, "\n")
	if len(newText) < 2 {
		return true, nil
	}
	date, DateErr := time.Parse("2006-01-02 15:04", newText[2])
	if DateErr != nil {
		return true, nil
	}

	wk := &Workout{Title: newText[0], Description: newText[1], CreatedAt: date, Status: 1}
	w, err := s.CreateWorkout(wk)
	if err != nil {
		return true, nil
	}
	return false, w
}

func (s WorkoutRepository) CallbackDeleteWorkout(update tgbotapi.Update) bool {
	id, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		return false
	}

	workout, err := s.ReadWorkout(id)
	if err != nil {
		return false
	}

	_, err = s.DeleteWorkout(workout)
	if err != nil {
		return false
	}

	return true
}

func (s WorkoutRepository) RegisterUserWorkout(user *users.User, workoutId int) *WorkoutUser {
	getWorkoutUser, err := s.FindBy(user.Id, workoutId)

	if err != nil {
		newWorkoutUser := WorkoutUser{
			UserID:    user.Id,
			WorkoutId: uint(workoutId),
		}
		getWorkoutUser, err = s.CreateWorkoutUser(&newWorkoutUser)
		if err != nil {
			panic(err)
		}
	}

	return getWorkoutUser
}

func (s WorkoutRepository) LeaveUserWorkout(user *users.User, workoutId int) {
	getWorkoutUser, err := s.FindBy(user.Id, workoutId)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = s.Delete(getWorkoutUser)
	if err != nil {
		fmt.Println(err)
		return
	}
}

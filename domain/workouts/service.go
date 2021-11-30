package workouts

// WorkoutService defines author service behavior.
type WorkoutService interface {
	CreateWorkout(*Workout) (*Workout, error)
	ReadWorkout(id int) (*Workout, error)
	UpdateWorkout(*Workout) (*Workout, error)
	ListWorkouts() ([]Workout, error)
	ListWorkoutMembers(workoutId int) ([]WorkoutUser, error)
	CreateWorkoutUser(user *WorkoutUser) (*WorkoutUser, error)
	FindBy(userId int, workoutId int) (*WorkoutUser, error)
	Delete(workoutUser *WorkoutUser) (bool, error)
}

// Service struct handles author business logic tasks.
type Service struct {
	repository WorkoutRepository
}

func (svc *Service) CreateWorkout(Workout *Workout) (*Workout, error) {
	return svc.repository.CreateWorkout(Workout)
}

func (svc *Service) ReadWorkout(id int) (*Workout, error) {
	return svc.repository.ReadWorkout(id)
}
func (svc *Service) UpdateWorkout(Workout *Workout) (*Workout, error) {
	return svc.repository.UpdateWorkout(Workout)
}

func (svc *Service) ListWorkouts() ([]Workout, error) {
	return svc.repository.ListWorkouts()
}

func (svc *Service) ListWorkoutMembers(workoutId int) ([]WorkoutUser, error) {
	return svc.repository.ListWorkoutMembers(workoutId)
}

func (svc *Service) CreateWorkoutUser(WorkoutUser *WorkoutUser) (*WorkoutUser, error) {
	return svc.repository.CreateWorkoutUser(WorkoutUser)
}

func (svc *Service) FindBy(userId int, workoutId int) (*WorkoutUser, error) {
	return svc.repository.FindBy(userId, workoutId)
}

func (svc *Service) Delete(workoutUser *WorkoutUser) (bool, error) {
	return svc.repository.Delete(workoutUser)
}

// NewService creates a new service struct
func NewService(repository WorkoutRepository) *Service {
	return &Service{repository: repository}
}

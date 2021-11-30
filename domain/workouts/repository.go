package workouts

type WorkoutRepository interface {
	CreateWorkout(*Workout) (*Workout, error)
	UpdateWorkout(*Workout) (*Workout, error)
	ReadWorkout(int) (*Workout, error)
	ListWorkouts() ([]Workout, error)
	ListWorkoutMembers(workoutId int) ([]WorkoutUser, error)
	CreateWorkoutUser(*WorkoutUser) (*WorkoutUser, error)
	FindBy(userId int, workoutId int) (*WorkoutUser, error)
	Delete(workoutUser *WorkoutUser) (bool, error)
}

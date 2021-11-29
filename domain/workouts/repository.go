package workouts

type WorkoutRepository interface {
	CreateWorkout(*Workout) (*Workout, error)
	UpdateWorkout(*Workout) (*Workout, error)
	ReadWorkout(int) (*Workout, error)
	ListWorkouts() ([]Workout, error)
	CreateWorkoutUser(*WorkoutUser) (*WorkoutUser, error)
}

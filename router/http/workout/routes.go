package workout

import (
	"github.com/gin-gonic/gin"
	"telegramStravaBot/domain/workouts"
)

// NewRoutesFactory create and returns a factory to create routes for the workout
func NewRoutesFactory(group *gin.RouterGroup) func(service workouts.WorkoutService) {
	RoutesFactory := func(service workouts.WorkoutService) {
		controller := new(Controller)
		controller.Service = service
		group.GET("/", controller.Index)
		group.GET("/:id", controller.Show)
		group.POST("/", controller.CreateWorkout)
	}

	return RoutesFactory
}

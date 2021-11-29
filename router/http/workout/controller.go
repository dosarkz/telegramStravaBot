package workout

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"telegramStravaBot/domain/workouts"
)

type Controller struct {
	Service workouts.WorkoutService
}

func (u *Controller) Index(c *gin.Context) {
	results, err := u.Service.ListWorkouts()
	if err != nil {
		c.JSON(404, gin.H{"message": "Workouts Not found"})
		return
	}

	var responseItems = make([]workouts.WorkoutResponse, len(results))

	for i, element := range results {
		responseItems[i] = *workouts.ToResponseModel(&element)
	}

	response := &workouts.ListWorkoutResponse{
		Data: responseItems,
	}

	c.JSON(http.StatusOK, response)
}

func (u *Controller) Create(c *gin.Context) {

	c.JSON(201, gin.H{"message": "Workout created"})
}

func (u *Controller) Show(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	result, err := u.Service.ReadWorkout(i)
	if err != nil {
		c.JSON(404, gin.H{"message": "Not found"})
		return
	}

	c.JSON(http.StatusOK, *workouts.ToResponseModel(result))
}

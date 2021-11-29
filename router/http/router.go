package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"telegramStravaBot/domain/workouts"
	"telegramStravaBot/router/http/middlewares"
	"telegramStravaBot/router/http/workout"
	"telegramStravaBot/utils/services"
)

// NewHTTPHandler returns the HTTP requests handler
func NewHTTPHandler(workoutService workouts.WorkoutService) http.Handler {
	router := gin.Default()
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))
	router.Use(services.JSONMiddleware())
	router.Use(middlewares.Authenticate())

	api := router.Group("/api/v1")

	workoutGroup := api.Group("/workouts")
	workout.NewRoutesFactory(workoutGroup)(workoutService)

	router.NoRoute(func(c *gin.Context) {
		// In gin this is how you return a JSON response
		c.JSON(404, gin.H{"message": "Not found"})
	})

	return router
}

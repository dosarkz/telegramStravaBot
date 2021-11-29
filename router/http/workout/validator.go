package workout

type WorkoutValidator struct {
	Title       string `binding:"required" json:"title"`
	Description string `binding:"required" json:"description"`
}

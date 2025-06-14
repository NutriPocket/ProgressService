package model

type ExerciseDTO struct {
	UserID         string  `json:"userId" binding:"required"`
	ExerciseName   string  `json:"exerciseName" binding:"required"`
	CaloriesBurned float64 `json:"caloriesBurned" binding:"required,gt=0"`
}

// ExerciseData represents an exercise performed by a user on a specific day
type ExerciseData struct {
	ID        uint64 `json:"id"`
	CreatedAt string `json:"createdAt"`
	ExerciseDTO
}

type AllExercisesInDay struct {
	TotalBurned float64        `json:"totalBurned"`
	Exercises   []ExerciseData `json:"exercises"`
}

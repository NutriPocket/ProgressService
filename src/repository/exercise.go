// Package repository provides structs and methods to interact with the database.
package repository

import (
	"fmt"

	"github.com/NutriPocket/ProgressService/database"
	"github.com/NutriPocket/ProgressService/model"
)

// IExerciseRepository is an interface that contains the methods that will implement a repository struct that interact with the exercise_by_day table.
type IExerciseRepository interface {
	CreateExercise(data *model.ExerciseDTO) (model.ExerciseData, error)
	GetExerciseById(id uint64, data *model.ExerciseData) error
	GetExercisesByUserIdAndDate(userId string, date string) (model.AllExercisesInDay, error)
	UpdateExercise(id uint64, data *model.ExerciseDTO) (model.ExerciseData, error)
	DeleteExercise(id uint64) error
}

type ExerciseRepository struct {
	db IDatabase
}

func NewExerciseRepository(db IDatabase) (*ExerciseRepository, error) {
	var err error

	if db == nil {
		db, err = database.GetPoolConnection()
		if err != nil {
			log.Errorf("Failed to connect to database")
			return nil, err
		}
	}

	return &ExerciseRepository{
		db: db,
	}, nil
}

func (r *ExerciseRepository) CreateExercise(data *model.ExerciseDTO) (model.ExerciseData, error) {
	res := r.db.Exec(`
        INSERT INTO exercise_by_day (user_id, exercise_name, calories_burned)
        VALUES (?, ?, ?);
    `,
		data.UserID, data.ExerciseName, data.CaloriesBurned,
	)

	if res.Error != nil {
		log.Errorf("Failed to create exercise for user %s: %v", data.UserID, res.Error)
		return model.ExerciseData{}, res.Error
	}

	// Get the last inserted ID
	var lastID uint64
	idRes := r.db.Raw("SELECT LAST_INSERT_ID()").Scan(&lastID)
	if idRes.Error != nil {
		log.Errorf("Failed to get last inserted ID: %v", idRes.Error)
		return model.ExerciseData{}, idRes.Error
	}

	// Retrieve the created exercise
	var createdExercise model.ExerciseData
	err := r.GetExerciseById(lastID, &createdExercise)
	return createdExercise, err
}

func (r *ExerciseRepository) GetExerciseById(id uint64, data *model.ExerciseData) error {
	res := r.db.Raw(`
        SELECT id, user_id, exercise_name, calories_burned, created_at
        FROM exercise_by_day
        WHERE id = ?
        LIMIT 1;
    `,
		id,
	).Scan(data)

	if res.Error != nil {
		return res.Error
	}

	if data.ID == 0 {
		return &model.NotFoundError{
			Title:  "Exercise not found",
			Detail: "No exercise found with ID " + fmt.Sprintf("%d", id),
		}
	}

	return nil
}

func (r *ExerciseRepository) UpdateExercise(id uint64, data *model.ExerciseDTO) (model.ExerciseData, error) {
	// Update the exercise
	res := r.db.Exec(`
		UPDATE exercise_by_day
		SET exercise_name = ?, calories_burned = ?
		WHERE id = ?;
	`,
		data.ExerciseName, data.CaloriesBurned, id,
	)

	if res.Error != nil {
		log.Errorf("Failed to update exercise with ID %d: %v", id, res.Error)
		return model.ExerciseData{}, res.Error
	}

	// Retrieve the updated exercise
	var updatedExercise model.ExerciseData
	err := r.GetExerciseById(id, &updatedExercise)
	return updatedExercise, err
}

func (r *ExerciseRepository) DeleteExercise(id uint64) error {
	// Delete the exercise
	res := r.db.Exec(`
        DELETE FROM exercise_by_day
        WHERE id = ?;
    `,
		id,
	)

	if res.Error != nil {
		log.Errorf("Failed to delete exercise with ID %d: %v", id, res.Error)
		return res.Error
	}

	return nil
}

func (r *ExerciseRepository) GetExercisesByUserIdAndDate(userId string, date string) (model.AllExercisesInDay, error) {
	// First, get all exercises for the day
	var exercises []model.ExerciseData
	res := r.db.Raw(`
        SELECT id, user_id, exercise_name, calories_burned, created_at
        FROM exercise_by_day
        WHERE user_id = ?
        AND DATE(created_at) = ?
        ORDER BY created_at ASC;
    `,
		userId, date,
	).Scan(&exercises)

	if res.Error != nil {
		log.Errorf("Failed to get exercises for user %s on date %s: %v", userId, date, res.Error)
		return model.AllExercisesInDay{}, res.Error
	}

	// Then, calculate the total calories burned
	var totalBurned float64
	sumRes := r.db.Raw(`
        SELECT COALESCE(SUM(calories_burned), 0) as total
        FROM exercise_by_day
        WHERE user_id = ?
        AND DATE(created_at) = ?;
    `,
		userId, date,
	).Scan(&totalBurned)

	if sumRes.Error != nil {
		log.Errorf("Failed to calculate total calories burned for user %s on date %s: %v", userId, date, sumRes.Error)
		return model.AllExercisesInDay{}, sumRes.Error
	}

	// Construct the result object
	result := model.AllExercisesInDay{
		TotalBurned: totalBurned,
		Exercises:   exercises,
	}

	return result, nil
}

package service

import (
	"time"

	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/repository"
)

// IExerciseService defines the interface for exercise-related operations
type IExerciseService interface {
	CreateExercise(data *model.ExerciseDTO) (model.ExerciseData, error)
	GetExercisesByUserIdAndDate(userId string, date string) (model.AllExercisesInDay, error)
	UpdateExercise(id uint64, userId string, data *model.ExerciseDTO) (model.ExerciseData, error)
	DeleteExercise(id uint64, userId string) error
}

// ExerciseService implements the IExerciseService interface
type ExerciseService struct {
	r repository.IExerciseRepository
}

// NewExerciseService creates a new ExerciseService instance
func NewExerciseService(r repository.IExerciseRepository) (*ExerciseService, error) {
	var err error

	if r == nil {
		r, err = repository.NewExerciseRepository(nil)
		if err != nil {
			return nil, err
		}
	}

	return &ExerciseService{
		r: r,
	}, nil
}

// CreateExercise adds a new exercise record
func (s *ExerciseService) CreateExercise(data *model.ExerciseDTO) (model.ExerciseData, error) {
	exercise, err := s.r.CreateExercise(data)
	if err != nil {
		log.Errorf("Failed to create exercise: %v", err)
		return model.ExerciseData{}, err
	}

	return exercise, nil
}

// GetExercisesByUserIdAndDate retrieves exercises for a user on a specific date
func (s *ExerciseService) GetExercisesByUserIdAndDate(userId string, date string) (model.AllExercisesInDay, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02") // Default to today if no date is provided
	} else {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			return model.AllExercisesInDay{}, &model.ValidationError{
				Title:  "Invalid date format",
				Detail: "Date must be in YYYY-MM-DD format",
			}
		}
	}

	exercises, err := s.r.GetExercisesByUserIdAndDate(userId, date)
	if err != nil {
		log.Errorf("Failed to get exercises for user %s on date %s: %v", userId, date, err)
		return model.AllExercisesInDay{}, err
	}

	return exercises, nil
}

// UpdateExercise updates an existing exercise
func (s *ExerciseService) UpdateExercise(id uint64, userId string, data *model.ExerciseDTO) (model.ExerciseData, error) {
	var existingExercise model.ExerciseData
	err := s.r.GetExerciseById(id, &existingExercise)
	if err != nil {
		return model.ExerciseData{}, err
	}

	if existingExercise.UserID != userId {
		return model.ExerciseData{}, &model.AuthenticationError{
			Title:  "Unauthorized",
			Detail: "You are not authorized to update this exercise",
		}
	}

	exercise, err := s.r.UpdateExercise(id, data)
	if err != nil {
		log.Errorf("Failed to update exercise with ID %d: %v", id, err)
		return model.ExerciseData{}, err
	}

	return exercise, nil
}

// DeleteExercise removes an exercise by its ID
func (s *ExerciseService) DeleteExercise(id uint64, userId string) error {
	var existingExercise model.ExerciseData
	err := s.r.GetExerciseById(id, &existingExercise)
	if err != nil {
		return err
	}

	if existingExercise.UserID != userId {
		return &model.AuthenticationError{
			Title:  "Unauthorized",
			Detail: "You are not authorized to delete this exercise",
		}
	}

	err = s.r.DeleteExercise(id)
	if err != nil {
		log.Errorf("Failed to delete exercise with ID %d: %v", id, err)
		return err
	}

	return nil
}

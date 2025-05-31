// Package repository provides structs and methods to interact with the database.
package repository

import (
	"github.com/NutriPocket/ProgressService/database"
	"github.com/NutriPocket/ProgressService/model"
	"github.com/go-sql-driver/mysql"
)

// IRoutineRepository is an interface that contains the methods that will implement a repository struct that interact with the users table.
type IRoutineRepository interface {
	CreateRoutine(data *model.RoutineDTO) (model.RoutineData, error)
	GetRoutinesByUserId(userId string, data *[]model.RoutineData) error
	GetRoutineBySchedule(userId string, schedule *model.Schedule) (model.RoutineData, error)
	GetRoutinesByInterval(userId string, schedule *model.Schedule) ([]model.RoutineData, error)
	DeleteRoutineBySchedule(userId string, schedule *model.Schedule) error
}

type RoutineRepository struct {
	db IDatabase
}

func NewRoutineRepository(db IDatabase) (*RoutineRepository, error) {
	var err error

	if db == nil {
		db, err = database.GetPoolConnection()
		if err != nil {
			log.Errorf("Failed to connect to database")
			return nil, err
		}
	}

	return &RoutineRepository{
		db: db,
	}, nil
}

func (r *RoutineRepository) CreateRoutine(data *model.RoutineDTO) (model.RoutineData, error) {
	res := r.db.Exec(`
		INSERT INTO user_routines (user_id, name, description, day, start_hour, end_hour)
		VALUES (?, ?, ?, ?, ?, ?);
	`,
		data.UserID, data.Name, data.Description, data.Day, data.StartHour, data.EndHour,
	)

	if res.Error != nil {
		// Check for MySQL duplicate entry error (error code 1062)
		log.Errorf("Failed to create a routine for user %s: %v", data.UserID, res.Error)

		if mysqlErr, ok := res.Error.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return model.RoutineData{}, &model.ConflictError{
				Title:  "Routine already exists",
				Detail: "A routine with the same schedule already exists for this user",
			}
		}

		return model.RoutineData{}, res.Error
	}

	ret, err := r.GetRoutineBySchedule(
		data.UserID,
		&model.Schedule{
			Day:       data.Day,
			StartHour: data.StartHour,
			EndHour:   data.EndHour,
		})

	return ret, err
}

func (r *RoutineRepository) GetRoutinesByInterval(userId string, schedule *model.Schedule) ([]model.RoutineData, error) {
	var routines []model.RoutineData

	res := r.db.Raw(`
		SELECT user_id, name, description, day, start_hour, end_hour, created_at, updated_at
		FROM user_routines
		WHERE day = ? AND end_hour >= ? AND start_hour < ? AND user_id = ?;
	`,
		schedule.Day, schedule.StartHour, schedule.EndHour, userId,
	).Scan(&routines)

	if res.Error != nil {
		return nil, res.Error
	}

	return routines, nil
}

func (r *RoutineRepository) GetRoutineBySchedule(userId string, schedule *model.Schedule) (model.RoutineData, error) {
	var routine model.RoutineData

	res := r.db.Raw(`
		SELECT user_id, name, description, day, start_hour, end_hour, created_at, updated_at
		FROM user_routines
		WHERE day = ? AND start_hour = ? AND end_hour = ? AND user_id = ?
		LIMIT 1;`,
		schedule.Day, schedule.StartHour, schedule.EndHour, userId,
	).Scan(&routine)

	if res.Error != nil {
		return model.RoutineData{}, res.Error
	}

	if routine.UserID == "" {
		return model.RoutineData{}, &model.NotFoundError{
			Title:  "Routine not found",
			Detail: "No routine found for the given schedule",
		}
	}

	return routine, nil
}

func (r *RoutineRepository) GetRoutinesByUserId(userId string, data *[]model.RoutineData) error {
	res := r.db.Raw(`
		SELECT user_id, name, description, day, start_hour, end_hour, created_at, updated_at
		FROM user_routines
		WHERE user_id = ?
	`,
		userId,
	).Scan(data)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *RoutineRepository) DeleteRoutineBySchedule(userId string, schedule *model.Schedule) error {
	res := r.db.Exec(`
		DELETE FROM user_routines
		WHERE day = ? AND start_hour = ? AND end_hour = ? AND user_id = ?;
	`,
		schedule.Day, schedule.StartHour, schedule.EndHour, userId,
	)

	if res.Error != nil {
		log.Errorf(
			"Failed to delete routine for schedule %s %d-%d for user %s: %v",
			schedule.Day, schedule.StartHour, schedule.EndHour, userId, res.Error,
		)
		return res.Error
	}

	return nil
}

// Package repository provides structs and methods to interact with the database.
package repository

import (
	"github.com/NutriPocket/ProgressService/database"
	"github.com/NutriPocket/ProgressService/model"
)

// IObjectiveRepository is an interface that contains the methods that will implement a repository struct that interact with the users table.
type IObjectiveRepository interface {
	CreateObjective(data *model.ObjectiveData) (model.ObjectiveData, error)
	ReplaceObjective(data *model.ObjectiveData) (model.ObjectiveData, error)
	GetObjectiveByUserId(userId string, data *model.ObjectiveData) error
}

type ObjectiveRepository struct {
	db IDatabase
}

func NewObjectiveRepository(db IDatabase) (*ObjectiveRepository, error) {
	var err error

	if db == nil {
		db, err = database.GetPoolConnection()
		if err != nil {
			log.Errorf("Failed to connect to database")
			return nil, err
		}
	}

	return &ObjectiveRepository{
		db: db,
	}, nil
}

func (r *ObjectiveRepository) CreateObjective(data *model.ObjectiveData) (model.ObjectiveData, error) {
	res := r.db.Exec(`
		INSERT INTO objective (user_id, weight, muscle_mass, fat_mass, bone_mass, deadline)
		VALUES (?, ?, ?, ?, ?, ?);
	`,
		data.UserID, data.Weight, data.MuscleMass, data.FatMass, data.BoneMass, data.Deadline,
	)

	if res.Error != nil {
		log.Errorf("Failed to create anthropometric data for user %s: %v", data.UserID, res.Error)
		return model.ObjectiveData{}, res.Error
	}

	var ret model.ObjectiveData
	err := r.GetObjectiveByUserId(data.UserID, &ret)

	return ret, err
}

func (r *ObjectiveRepository) ReplaceObjective(data *model.ObjectiveData) (model.ObjectiveData, error) {
	res := r.db.Exec(`
		UPDATE objective
		SET weight = ?, muscle_mass = ?, fat_mass = ?, bone_mass = ?, deadline = ?
		WHERE user_id = ?;
	`,
		data.Weight, data.MuscleMass, data.FatMass, data.BoneMass, data.Deadline, data.UserID,
	)

	if res.Error != nil {
		log.Errorf("Failed to update anthropometric data for user %s: %v", data.UserID, res.Error)
		return model.ObjectiveData{}, res.Error
	}

	var ret model.ObjectiveData
	err := r.GetObjectiveByUserId(data.UserID, &ret)

	return ret, err
}

func (r *ObjectiveRepository) GetObjectiveByUserId(userId string, data *model.ObjectiveData) error {
	res := r.db.Raw(`
		SELECT user_id, weight, muscle_mass, fat_mass, bone_mass, created_at, deadline
		FROM objective
		WHERE user_id = ?
		LIMIT 1;`,
		userId,
	).Scan(data)

	if res.Error != nil {
		return res.Error
	}

	if data.UserID == "" {
		return &model.NotFoundError{
			Title:  "Objective data not found",
			Detail: "No objective data not found for user " + userId,
		}
	}

	return nil
}

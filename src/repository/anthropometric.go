// Package repository provides structs and methods to interact with the database.
package repository

import (
	"time"

	"github.com/NutriPocket/ProgressService/database"
	"github.com/NutriPocket/ProgressService/model"
)

// IAnthropometricRepository is an interface that contains the methods that will implement a repository struct that interact with the users table.
type IAnthropometricRepository interface {
	CreateData(data *model.AnthropometricData) (model.AnthropometricData, error)
	ReplaceTodayData(data *model.AnthropometricData) (model.AnthropometricData, error)
	GetDataByUserIdAndDate(userId string, date string, data *model.AnthropometricData) error
	GetTodayDataByUserId(userId string, data *model.AnthropometricData) error
	GetAllDataByUserId(userId string) ([]model.AnthropometricData, error)
}

type AnthropometricRepository struct {
	db IDatabase
}

func NewAnthropometricRepository(db IDatabase) (*AnthropometricRepository, error) {
	var err error

	if db == nil {
		db, err = database.GetPoolConnection()
		if err != nil {
			log.Errorf("Failed to connect to database")
			return nil, err
		}
	}

	return &AnthropometricRepository{
		db: db,
	}, nil
}

func (r *AnthropometricRepository) GetTodayDataByUserId(userId string, data *model.AnthropometricData) error {
	var date string
	return r.GetDataByUserIdAndDate(userId, date, data)
}

func (r *AnthropometricRepository) CreateData(data *model.AnthropometricData) (model.AnthropometricData, error) {
	res := r.db.Exec(`
		INSERT INTO anthropometric_data (user_id, weight, muscle_mass, fat_mass, bone_mass)
		VALUES (?, ?, ?, ?, ?);
	`,
		data.UserID, data.Weight, data.MuscleMass, data.FatMass, data.BoneMass,
	)

	if res.Error != nil {
		log.Errorf("Failed to create anthropometric data for user %s: %v", data.UserID, res.Error)
		return model.AnthropometricData{}, res.Error
	}

	var ret model.AnthropometricData
	err := r.GetTodayDataByUserId(data.UserID, &ret)

	return ret, err
}

func (r *AnthropometricRepository) ReplaceTodayData(data *model.AnthropometricData) (model.AnthropometricData, error) {
	res := r.db.Exec(`
		UPDATE anthropometric_data
		SET weight = ?, muscle_mass = ?, fat_mass = ?, bone_mass = ?
		WHERE user_id = ? 
			AND DATE(created_at) = CURDATE();
	`,
		data.Weight, data.MuscleMass, data.FatMass, data.BoneMass, data.UserID,
	)

	if res.Error != nil {
		log.Errorf("Failed to update anthropometric data for user %s: %v", data.UserID, res.Error)
		return model.AnthropometricData{}, res.Error
	}

	var ret model.AnthropometricData
	err := r.GetTodayDataByUserId(data.UserID, &ret)

	return ret, err
}

func (r *AnthropometricRepository) GetDataByUserIdAndDate(userId string, date string, data *model.AnthropometricData) error {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	res := r.db.Raw(`
		SELECT user_id, weight, muscle_mass, fat_mass, bone_mass, created_at
		FROM anthropometric_data 
		WHERE user_id = ? 
			AND DATE(created_at) = ?
		LIMIT 1;`,
		userId, date,
	).Scan(&data)

	if res.Error != nil {
		return res.Error
	}

	if data.UserID == "" {
		return &model.NotFoundError{
			Title:  "Anthropometric data not found",
			Detail: "No anthropometric data not found for user " + userId + " on date " + date,
		}
	}

	return nil
}

func (r *AnthropometricRepository) GetAllDataByUserId(userId string) ([]model.AnthropometricData, error) {
	var data []model.AnthropometricData = make([]model.AnthropometricData, 0)

	res := r.db.Raw(`
		SELECT user_id, weight, muscle_mass, fat_mass, bone_mass, created_at
		FROM anthropometric_data 
		WHERE user_id = ? 
		ORDER BY created_at DESC;
	`, userId,
	).Scan(&data)

	if res.Error != nil {
		return nil, res.Error
	}

	return data, nil
}

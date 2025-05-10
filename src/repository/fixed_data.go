// Package repository provides structs and methods to interact with the database.
package repository

import (
	"errors"

	"github.com/NutriPocket/ProgressService/database"
	"github.com/NutriPocket/ProgressService/model"
	"github.com/go-sql-driver/mysql"
)

// IFixedDataRepository is an interface that contains the methods that will implement a repository struct that interact with the users table.
type IFixedDataRepository interface {
	CreateData(data *model.BaseFixedUserData) (model.FixedUserData, error)
	ReplaceData(data *model.BaseFixedUserData) (model.FixedUserData, error)
	GetBaseFixedUserData(userId string, data *model.BaseFixedUserData) error
	GetUserData(userId string, data *model.FixedUserData) error
}

type FixedDataRepository struct {
	db IDatabase
}

func NewFixedDataRepository(db IDatabase) (*FixedDataRepository, error) {
	var err error

	if db == nil {
		db, err = database.GetPoolConnection()
		if err != nil {
			log.Errorf("Failed to connect to database")
			return nil, err
		}
	}

	return &FixedDataRepository{
		db: db,
	}, nil
}

func (r *FixedDataRepository) CreateData(data *model.BaseFixedUserData) (model.FixedUserData, error) {
	res := r.db.Exec(`
		INSERT INTO fixed_user_data (user_id, height, birthday)
		VALUES (?, ?, ?);
	`,
		data.UserID, data.Height, data.Birthday,
	)

	if res.Error != nil {
		log.Errorf("Failed to create fixed user data for user %s: %v", data.UserID, res.Error)

		if errors.Is(res.Error, &mysql.MySQLError{Number: 1062}) {
			return model.FixedUserData{}, &model.EntityAlreadyExistsError{
				Detail: "User fixed data already exists for the user " + data.UserID,
				Title:  "User fixed data already exists",
			}
		}

		return model.FixedUserData{}, res.Error
	}

	var ret model.FixedUserData
	err := r.GetUserData(data.UserID, &ret)
	return ret, err
}

func (r *FixedDataRepository) ReplaceData(data *model.BaseFixedUserData) (model.FixedUserData, error) {
	res := r.db.Exec(`
		UPDATE fixed_user_data
		SET height = ?, birthday = ?
		WHERE user_id = ?
	`,
		data.Height, data.Birthday, data.UserID,
	)

	if res.Error != nil {
		log.Errorf("Failed to update fixed user data for user %s: %v", data.UserID, res.Error)
		return model.FixedUserData{}, res.Error
	}

	var ret model.FixedUserData
	err := r.GetUserData(data.UserID, &ret)
	return ret, err
}

func (r *FixedDataRepository) GetBaseFixedUserData(userId string, data *model.BaseFixedUserData) error {
	res := r.db.Raw(`
		SELECT user_id, height, birthday 
		FROM fixed_user_data 
		WHERE user_id = ?
		LIMIT 1`,
		userId,
	).Scan(&data)

	log.Warningf("Query res: %v", data)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *FixedDataRepository) GetUserData(userId string, data *model.FixedUserData) error {
	res := r.db.Raw(`
		SELECT user_id, height, FLOOR(DATEDIFF(CURRENT_DATE(), birthday) / 365.25) AS age 
		FROM fixed_user_data 
		WHERE user_id = ?
		LIMIT 1`,
		userId,
	).Scan(data)

	log.Warningf("Query res: %v", data)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

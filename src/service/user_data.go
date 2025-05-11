package service

import (
	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/repository"
)

type IUserDataService interface {
	PutAnthropometricData(data *model.AnthropometricData) (model.AnthropometricData, error, bool)
	GetAnthropometricDataByUserAndDay(userId string, date string) (model.AnthropometricData, error)
	GetAllAnthropometricDataByUser(userId string) ([]model.AnthropometricData, error)
	PutFixedData(data *model.BaseFixedUserData) (model.FixedUserData, error, bool)
	GetFixedDataByUser(userId string) (model.FixedUserData, error)
	GetBaseFixedUserDataByUser(userId string) (model.BaseFixedUserData, error)
}

type UserDataService struct {
	ar  repository.IAnthropometricRepository
	fdr repository.IFixedDataRepository
}

func NewUserDataService(ar repository.IAnthropometricRepository, fdr repository.IFixedDataRepository) (*UserDataService, error) {
	var err error

	if ar == nil {
		ar, err = repository.NewAnthropometricRepository(nil)
		if err != nil {
			return nil, err
		}
	}

	if fdr == nil {
		fdr, err = repository.NewFixedDataRepository(nil)
		if err != nil {
			return nil, err
		}
	}

	return &UserDataService{
		ar:  ar,
		fdr: fdr,
	}, nil
}

func (s *UserDataService) PutAnthropometricData(data *model.AnthropometricData) (ret model.AnthropometricData, err error, created bool) {
	var storedData *model.AnthropometricData = &model.AnthropometricData{}
	err = s.ar.GetTodayDataByUserId(data.UserID, storedData)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); ok {
			ret, err = s.ar.CreateData(data)
			created = true
			return
		}

		log.Errorf("Failed to check current anthropometric data for user %s: %v", data.UserID, err)
		return
	}

	storedData.Weight = data.Weight
	if data.MuscleMass != nil {
		storedData.MuscleMass = data.MuscleMass
	}

	if data.FatMass != nil {
		storedData.FatMass = data.FatMass
	}

	if data.BoneMass != nil {
		storedData.BoneMass = data.BoneMass
	}

	ret, err = s.ar.ReplaceTodayData(storedData)
	return
}

func (s *UserDataService) GetAnthropometricDataByUserAndDay(userId string, date string) (model.AnthropometricData, error) {
	var ret model.AnthropometricData
	err := s.ar.GetDataByUserIdAndDate(userId, date, &ret)

	return ret, err
}

func (s *UserDataService) GetAllAnthropometricDataByUser(userId string) ([]model.AnthropometricData, error) {
	return s.ar.GetAllDataByUserId(userId)
}

func (s *UserDataService) PutFixedData(data *model.BaseFixedUserData) (ret model.FixedUserData, err error, created bool) {
	var storedData *model.BaseFixedUserData = &model.BaseFixedUserData{}
	err = s.fdr.GetBaseFixedUserData(data.UserID, storedData)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); ok {
			ret, err = s.fdr.CreateData(data)
		created = true
		return
		}

		log.Errorf("Failed to get fixed user data for user %s: %v", data.UserID, err)
		return
	}

	if len(data.Birthday) != 0 {
		storedData.Birthday = data.Birthday
	}

	if data.Height != 0 {
		storedData.Height = data.Height
	}

	ret, err = s.fdr.ReplaceData(storedData)
	return
}

func (s *UserDataService) GetFixedDataByUser(userId string) (model.FixedUserData, error) {
	var ret model.FixedUserData
	err := s.fdr.GetUserData(userId, &ret)

	return ret, err
}

func (s *UserDataService) GetBaseFixedUserDataByUser(userId string) (model.BaseFixedUserData, error) {
	var ret model.BaseFixedUserData
	err := s.fdr.GetBaseFixedUserData(userId, &ret)

	return ret, err
}

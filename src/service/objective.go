package service

import (
	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/repository"
)

type IObjectiveService interface {
	PutObjective(data *model.ObjectiveData) (model.ObjectiveData, error, bool)
	GetObjectiveByUser(userId string) (model.ObjectiveData, error)
}

type ObjectiveService struct {
	r repository.IObjectiveRepository
}

func NewObjectiveService(r repository.IObjectiveRepository) (*ObjectiveService, error) {
	var err error

	if r == nil {
		r, err = repository.NewObjectiveRepository(nil)
		if err != nil {
			return nil, err
		}
	}

	return &ObjectiveService{
		r: r,
	}, nil
}

func (s *ObjectiveService) PutObjective(data *model.ObjectiveData) (ret model.ObjectiveData, err error, created bool) {
	var storedData *model.ObjectiveData = &model.ObjectiveData{}
	err = s.r.GetObjectiveByUserId(data.UserID, storedData)
	if err != nil {
		log.Errorf("Failed to check current anthropometric data for user %s: %v", data.UserID, err)
		return
	}

	if storedData.UserID == "" {
		ret, err = s.r.CreateObjective(data)
		created = true
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

	if data.Deadline != "" {
		storedData.Deadline = data.Deadline
	}

	ret, err = s.r.ReplaceObjective(storedData)
	return
}

func (s *ObjectiveService) GetObjectiveByUser(userId string) (model.ObjectiveData, error) {
	var ret model.ObjectiveData
	err := s.r.GetObjectiveByUserId(userId, &ret)

	return ret, err
}

package service

import (
	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/repository"
)

type IRoutineService interface {
	CreateRoutine(data *model.RoutineDTO) (model.RoutineData, error)
	GetRoutinesByUser(userId string, data *[]model.RoutineData) error
	GetFreeSchedules(users []string) (model.FreeSchedule, error)
	DeleteRutineBySchedule(userId string, schedule *model.Schedule) ([]model.RoutineData, error)
}

type RoutineService struct {
	r repository.IRoutineRepository
}

func NewRoutineService(r repository.IRoutineRepository) (*RoutineService, error) {
	var err error

	if r == nil {
		r, err = repository.NewRoutineRepository(nil)
		if err != nil {
			return nil, err
		}
	}

	return &RoutineService{
		r: r,
	}, nil
}

func (s *RoutineService) CreateRoutine(data *model.RoutineDTO) (model.RoutineData, error) {
	existentRoutines, err := s.r.GetRoutinesByInterval(
		data.UserID,
		&model.Schedule{
			Day:       data.Day,
			StartHour: data.StartHour,
			EndHour:   data.EndHour,
		})

	if err != nil {
		return model.RoutineData{}, err
	}

	if len(existentRoutines) > 0 {
		return model.RoutineData{}, &model.ConflictError{
			Title:  "Routine conflict",
			Detail: "There is already a routine scheduled in the same time interval or subinterval",
		}
	}

	ret, err := s.r.CreateRoutine(data)
	if err != nil {
		return model.RoutineData{}, err
	}

	return ret, nil
}

func (s *RoutineService) GetRoutinesByUser(userId string, data *[]model.RoutineData) error {
	err := s.r.GetRoutinesByUserId(userId, data)

	if err != nil {
		return err
	}

	return nil
}

func (s *RoutineService) getFreeHours(users []string) (map[string][]bool, error) {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	freeSchedules := make(map[string][]bool, 0)
	for _, day := range days {
		freeSchedules[day] = make([]bool, 24) // Assuming 24 hours in a day
		for j := range freeSchedules[day] {
			freeSchedules[day][j] = true
		}
	}

	for _, user := range users {
		var routines []model.RoutineData

		err := s.r.GetRoutinesByUserId(user, &routines)
		if err != nil {
			return nil, err
		}

		for _, routine := range routines {
			for i := routine.StartHour; i < routine.EndHour; i++ {
				if i >= 0 && i < len(freeSchedules[routine.Day]) {
					freeSchedules[routine.Day][i] = false
				}
			}
		}
	}

	return freeSchedules, nil
}

func (s *RoutineService) freeHoursToSchedule(freeSchedules map[string][]bool) []model.Schedule {
	schedules := make([]model.Schedule, 0)

	for day, hours := range freeSchedules {
		start := -1
		end := -1

		for i, isFree := range hours {
			if isFree {
				if start == -1 {
					start = i
					end = i + 1
				}

				if i > end {
					end = i
				}

				if i != len(hours)-1 {
					continue
				}
			}

			if start != -1 {
				schedules = append(schedules, model.Schedule{
					Day:       day,
					StartHour: start,
					EndHour:   end,
				})
				start = -1
				end = -1
			}
		}
	}

	return schedules
}

func (s *RoutineService) GetFreeSchedules(users []string) (model.FreeSchedule, error) {
	if len(users) == 0 {
		return model.FreeSchedule{
			Schedules: []model.Schedule{},
		}, nil
	}

	freeSchedules, err := s.getFreeHours(users)
	if err != nil {
		return model.FreeSchedule{}, err
	}

	data := model.FreeSchedule{
		Schedules: s.freeHoursToSchedule(freeSchedules),
	}

	if len(data.Schedules) == 0 {
		return model.FreeSchedule{}, &model.NotFoundError{
			Title:  "No free schedules found",
			Detail: "No free schedules found for the provided users",
		}
	}

	return data, nil
}

func (s *RoutineService) DeleteRutineBySchedule(userId string, schedule *model.Schedule) ([]model.RoutineData, error) {
	err := s.r.DeleteRoutineBySchedule(userId, schedule)
	if err != nil {
		return nil, err
	}

	var data []model.RoutineData

	s.GetRoutinesByUser(userId, &data)

	return data, nil
}

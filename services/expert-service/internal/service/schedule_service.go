package service

import (
	"expert-service/internal/model"
	"expert-service/internal/repository"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ScheduleService interface {
	CreateSchedule(req *model.CreateScheduleRequest) (*model.Schedule, error)
	GetSchedulesByExpertID(expertID uuid.UUID) ([]*model.Schedule, error)
	GetScheduleByID(id uuid.UUID) (*model.Schedule, error)
	UpdateSchedule(id uuid.UUID, req *model.CreateScheduleRequest) error
	DeleteSchedule(id uuid.UUID) error
}

type scheduleService struct {
	scheduleRepo repository.ScheduleRepository
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository) ScheduleService {
	return &scheduleService{
		scheduleRepo: scheduleRepo,
	}
}

func (s *scheduleService) CreateSchedule(req *model.CreateScheduleRequest) (*model.Schedule, error) {
	schedule := &model.Schedule{
		ID:        uuid.New(),
		ExpertID:  strconv.Itoa(req.ExpertID),
		DayOfWeek: strconv.Itoa(req.DayOfWeek),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		IsActive:  time.Now(),
		CreatedAt: time.Now(),
	}

	if err := s.scheduleRepo.Create(schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return schedule, nil
}

func (s *scheduleService) GetSchedulesByExpertID(expertID uuid.UUID) ([]*model.Schedule, error) {
	expertIDInt, err := strconv.Atoi(expertID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid expert ID format: %w", err)
	}
	return s.scheduleRepo.GetByExpertID(expertIDInt)
}

func (s *scheduleService) GetScheduleByID(id uuid.UUID) (*model.Schedule, error) {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}
	if schedule == nil {
		return nil, fmt.Errorf("schedule not found")
	}
	return schedule, nil
}

func (s *scheduleService) UpdateSchedule(id uuid.UUID, req *model.CreateScheduleRequest) error {
	schedule := &model.Schedule{
		ID:        id,
		ExpertID:  strconv.Itoa(req.ExpertID),
		DayOfWeek: strconv.Itoa(req.DayOfWeek),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		IsActive:  time.Now(),
	}

	return s.scheduleRepo.Update(schedule)
}

func (s *scheduleService) DeleteSchedule(id uuid.UUID) error {
	return s.scheduleRepo.Delete(id)
}

package service

import (
	"expert-service/internal/model"
	"expert-service/internal/repository"
	"fmt"
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
		ExpertID:  req.ExpertID,
		DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Timezone:  req.Timezone,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.scheduleRepo.Create(schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return schedule, nil
}

func (s *scheduleService) GetSchedulesByExpertID(expertID uuid.UUID) ([]*model.Schedule, error) {
	return s.scheduleRepo.GetByExpertID(expertID)
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
		ExpertID:  req.ExpertID,
		DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Timezone:  req.Timezone,
		IsActive:  true,
		UpdatedAt: time.Now(),
	}

	return s.scheduleRepo.Update(id, schedule)
}

func (s *scheduleService) DeleteSchedule(id uuid.UUID) error {
	return s.scheduleRepo.Delete(id)
}


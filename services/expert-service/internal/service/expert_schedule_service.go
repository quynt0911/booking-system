package service

import (
	"expert-service/internal/model"
	"expert-service/internal/repository"

	"github.com/google/uuid"
)

type ExpertScheduleService interface {
	CreateSchedule(schedule *model.ExpertSchedule) error
	GetSchedulesByExpertID(expertID uuid.UUID) ([]*model.ExpertSchedule, error)
	UpdateSchedule(schedule *model.ExpertSchedule) error
	DeleteSchedule(id uuid.UUID) error
}

type expertScheduleService struct {
	repo repository.ExpertScheduleRepository
}

func NewExpertScheduleService(repo repository.ExpertScheduleRepository) ExpertScheduleService {
	return &expertScheduleService{repo: repo}
}

func (s *expertScheduleService) CreateSchedule(schedule *model.ExpertSchedule) error {
	return s.repo.Create(schedule)
}

func (s *expertScheduleService) GetSchedulesByExpertID(expertID uuid.UUID) ([]*model.ExpertSchedule, error) {
	return s.repo.GetByExpertID(expertID)
}

func (s *expertScheduleService) UpdateSchedule(schedule *model.ExpertSchedule) error {
	return s.repo.Update(schedule)
}

func (s *expertScheduleService) DeleteSchedule(id uuid.UUID) error {
	return s.repo.Delete(id)
}

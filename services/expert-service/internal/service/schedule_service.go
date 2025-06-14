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
	UpdateSchedule(id uuid.UUID, req *model.UpdateScheduleRequest) error
	DeleteSchedule(id uuid.UUID) error
	GetSchedules(req *model.GetSchedulesRequest) ([]*model.Schedule, error)
	CancelSchedule(id uuid.UUID) error
	ConfirmSchedule(id uuid.UUID) error
	CompleteSchedule(id uuid.UUID) error
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
		ID:             uuid.New(),
		ExpertID:       req.ExpertID,
		UserID:         req.UserID,
		AvailabilityID: req.AvailabilityID,
		DayOfWeek:      req.DayOfWeek,
		Date:           req.Date,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Title:          req.Title,
		Description:    req.Description,
		MeetingLink:    req.MeetingLink,
		Notes:          req.Notes,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.scheduleRepo.Create(schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return schedule, nil
}

func (s *scheduleService) GetSchedulesByExpertID(expertID uuid.UUID) ([]*model.Schedule, error) {
	return s.scheduleRepo.GetByExpertID(expertID.String())
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

func (s *scheduleService) UpdateSchedule(id uuid.UUID, req *model.UpdateScheduleRequest) error {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get schedule for update: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("schedule with ID %s not found", id.String())
	}

	// Apply updates from request
	if req.ExpertID != nil {
		schedule.ExpertID = *req.ExpertID
	}
	if req.UserID != nil {
		schedule.UserID = *req.UserID
	}
	if req.AvailabilityID != nil {
		schedule.AvailabilityID = *req.AvailabilityID
	}
	if req.DayOfWeek != nil {
		schedule.DayOfWeek = *req.DayOfWeek
	}
	if req.Date != nil {
		schedule.Date = *req.Date
	}
	if req.StartTime != nil {
		schedule.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		schedule.EndTime = *req.EndTime
	}
	if req.Status != nil {
		schedule.Status = *req.Status
	}
	if req.Title != nil {
		schedule.Title = *req.Title
	}
	if req.Description != nil {
		schedule.Description = *req.Description
	}
	if req.MeetingLink != nil {
		schedule.MeetingLink = *req.MeetingLink
	}
	if req.Notes != nil {
		schedule.Notes = *req.Notes
	}
	if req.IsActive != nil {
		schedule.IsActive = *req.IsActive
	}

	schedule.UpdatedAt = time.Now()

	return s.scheduleRepo.Update(schedule)
}

func (s *scheduleService) DeleteSchedule(id uuid.UUID) error {
	return s.scheduleRepo.Delete(id)
}

// GetSchedules retrieves schedules based on provided filters
func (s *scheduleService) GetSchedules(req *model.GetSchedulesRequest) ([]*model.Schedule, error) {
	return s.scheduleRepo.GetSchedules(req)
}

// CancelSchedule updates a schedule's status to cancelled
func (s *scheduleService) CancelSchedule(id uuid.UUID) error {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get schedule for cancellation: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("schedule with ID %s not found", id.String())
	}

	schedule.Status = "cancelled"
	schedule.IsActive = false // Deactivate cancelled schedules
	schedule.UpdatedAt = time.Now()

	return s.scheduleRepo.Update(schedule)
}

// ConfirmSchedule updates a schedule's status to confirmed
func (s *scheduleService) ConfirmSchedule(id uuid.UUID) error {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get schedule for confirmation: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("schedule with ID %s not found", id.String())
	}

	schedule.Status = "confirmed"
	schedule.UpdatedAt = time.Now()

	return s.scheduleRepo.Update(schedule)
}

// CompleteSchedule updates a schedule's status to completed
func (s *scheduleService) CompleteSchedule(id uuid.UUID) error {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get schedule for completion: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("schedule with ID %s not found", id.String())
	}

	schedule.Status = "completed"
	schedule.UpdatedAt = time.Now()

	return s.scheduleRepo.Update(schedule)
}

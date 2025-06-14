package service

import (
	"encoding/json"
	"expert-service/internal/cache"
	"expert-service/internal/model"
	"expert-service/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// isTimeInRange checks if a time string falls within a start and end time range
func (s *expertAvailabilityService) isTimeInRange(timeStr, startTimeStr, endTimeStr string) bool {
	t, _ := time.ParseInLocation("15:04", timeStr, time.Local)
	start, _ := time.ParseInLocation("15:04", startTimeStr, time.Local)
	end, _ := time.ParseInLocation("15:04", endTimeStr, time.Local)
	return (t.Equal(start) || t.After(start)) && t.Before(end)
}

type ExpertAvailabilityService interface {
	CheckAvailability(req *model.CheckAvailabilityRequest) (bool, error)
	CreateOffTime(req *model.CreateOffTimeRequest) (*model.OffTime, error)
	GetExpertOffTimes(expertID string) ([]*model.OffTime, error)
	DeleteOffTime(id string) error
	CreateAvailability(req *model.CreateAvailabilityRequest) (*model.Availability, error)
	GetAvailabilityByID(id string) (*model.Availability, error)
	UpdateAvailability(id string, req *model.UpdateAvailabilityRequest) (*model.Availability, error)
	DeleteAvailability(id string) error
	GetAvailabilities(expertID string, startDate, endDate time.Time, isBooked *bool) ([]*model.Availability, error)
	BookAvailability(id string) error
	CreateRecurringAvailability(req *model.CreateRecurringAvailabilityRequest) ([]*model.Availability, error)
}

type expertAvailabilityService struct {
	expertRepo   repository.ExpertRepository
	scheduleRepo repository.ScheduleRepository
	offTimeRepo  repository.OffTimeRepository
	cache        cache.AvailabilityCache
}

func NewExpertAvailabilityService(
	expertRepo repository.ExpertRepository,
	scheduleRepo repository.ScheduleRepository,
	offTimeRepo repository.OffTimeRepository,
	cache cache.AvailabilityCache,
) ExpertAvailabilityService {
	return &expertAvailabilityService{
		expertRepo:   expertRepo,
		scheduleRepo: scheduleRepo,
		offTimeRepo:  offTimeRepo,
		cache:        cache,
	}
}

// CheckAvailability kiểm tra chuyên gia có rảnh không tại thời điểm yêu cầu
func (s *expertAvailabilityService) CheckAvailability(req *model.CheckAvailabilityRequest) (bool, error) {
	expertUUID, err := uuid.Parse(req.ExpertID)
	if err != nil {
		return false, fmt.Errorf("invalid expert ID format: %v", err)
	}
	expert, err := s.expertRepo.GetByID(expertUUID)
	if err != nil {
		return false, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
	}
	if expert == nil {
		return false, fmt.Errorf("không tìm thấy chuyên gia với ID %s", req.ExpertID)
	}

	// Parse date and time
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return false, fmt.Errorf("định dạng ngày không hợp lệ")
	}

	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s", req.ExpertID, req.Date)
	if cached, err := s.cache.GetAvailability(cacheKey); err == nil && cached != nil {
		var isAvailable bool
		if err := json.Unmarshal(cached, &isAvailable); err == nil {
			return isAvailable, nil
		}
	}

	// Check if expert is on off-time
	offTimes, err := s.offTimeRepo.GetByExpertIDAndDateRange(expertUUID, date)
	if err != nil {
		return false, fmt.Errorf("không thể kiểm tra thời gian nghỉ: %v", err)
	}
	if len(offTimes) > 0 {
		data, _ := json.Marshal(false)
		if err := s.cache.SetAvailability(cacheKey, data); err != nil {
			return false, fmt.Errorf("failed to cache availability: %w", err)
		}
		return false, nil
	}

	// Get expert's schedule for the day
	dayOfWeek := int(date.Weekday())
	schedules, err := s.scheduleRepo.GetByExpertIDAndDay(req.ExpertID, dayOfWeek)
	if err != nil {
		return false, fmt.Errorf("không thể lấy lịch làm việc: %v", err)
	}

	// Check if the requested time falls within any schedule
	isAvailable := false
	for _, schedule := range schedules {
		if s.isTimeInRange(req.Time, schedule.StartTime, schedule.EndTime) {
			isAvailable = true
			break
		}
	}

	// Cache the result
	data, _ := json.Marshal(isAvailable)
	if err := s.cache.SetAvailability(cacheKey, data); err != nil {
		return false, fmt.Errorf("failed to cache availability: %w", err)
	}
	return isAvailable, nil
}

// CreateOffTime tạo thời gian nghỉ cho chuyên gia
func (s *expertAvailabilityService) CreateOffTime(req *model.CreateOffTimeRequest) (*model.OffTime, error) {
	expertUUID, err := uuid.Parse(req.ExpertID)
	if err != nil {
		return nil, fmt.Errorf("invalid expert ID format: %v", err)
	}
	expert, err := s.expertRepo.GetByID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
	}
	if expert == nil {
		return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %s", req.ExpertID)
	}

	startDateTime, err := time.Parse("2006-01-02T15:04:05Z", req.StartDateTime)
	if err != nil {
		return nil, fmt.Errorf("định dạng thời gian bắt đầu không hợp lệ")
	}

	endDateTime, err := time.Parse("2006-01-02T15:04:05Z", req.EndDateTime)
	if err != nil {
		return nil, fmt.Errorf("định dạng thời gian kết thúc không hợp lệ")
	}

	if endDateTime.Before(startDateTime) {
		return nil, fmt.Errorf("thời gian kết thúc phải sau thời gian bắt đầu")
	}

	offTime := &model.OffTime{
		ExpertID:      expertUUID,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		Reason:        req.Reason,
		IsRecurring:   req.IsRecurring,
	}

	err = s.offTimeRepo.Create(offTime)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo thời gian nghỉ: %v", err)
	}

	// Invalidate cache
	s.cache.InvalidateExpert(req.ExpertID)
	return offTime, nil
}

// GetExpertOffTimes lấy danh sách thời gian nghỉ của chuyên gia
func (s *expertAvailabilityService) GetExpertOffTimes(expertID string) ([]*model.OffTime, error) {
	expertUUID, err := uuid.Parse(expertID)
	if err != nil {
		return nil, fmt.Errorf("invalid expert ID format: %v", err)
	}
	expert, err := s.expertRepo.GetByID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
	}
	if expert == nil {
		return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %s", expertID)
	}

	offTimes, err := s.offTimeRepo.GetByExpertID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách thời gian nghỉ: %v", err)
	}

	return offTimes, nil
}

// DeleteOffTime xóa thời gian nghỉ
func (s *expertAvailabilityService) DeleteOffTime(id string) error {
	offTimeID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid off time ID format: %v", err)
	}

	err = s.offTimeRepo.Delete(offTimeID)
	if err != nil {
		return fmt.Errorf("không thể xóa thời gian nghỉ: %v", err)
	}
	return nil
}

// CreateAvailability creates a new availability slot
func (s *expertAvailabilityService) CreateAvailability(req *model.CreateAvailabilityRequest) (*model.Availability, error) {
	expertUUID, err := uuid.Parse(req.ExpertID)
	if err != nil {
		return nil, fmt.Errorf("invalid expert ID format: %v", err)
	}
	expert, err := s.expertRepo.GetByID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
	}
	if expert == nil {
		return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %s", req.ExpertID)
	}

	availability := &model.Availability{
		ID:        uuid.New(),
		ExpertID:  req.ExpertID,
		Date:      req.Date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		IsBooked:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Lưu vào Redis với key pattern: availability:{expert_id}:{date}
	key := fmt.Sprintf("availability:%s:%s", req.ExpertID, req.Date)
	data, err := json.Marshal(availability)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal availability: %w", err)
	}

	if err := s.cache.SetAvailability(key, data); err != nil {
		return nil, fmt.Errorf("failed to cache availability: %w", err)
	}

	return availability, nil
}

// GetAvailabilityByID retrieves an availability slot by ID
func (s *expertAvailabilityService) GetAvailabilityByID(id string) (*model.Availability, error) {
	// Tìm trong Redis với pattern: availability:*:*
	// TODO: Implement search by ID in Redis
	// For now, return error as this might need a different Redis structure
	return nil, fmt.Errorf("get availability by ID not implemented with Redis")
}

// UpdateAvailability updates an existing availability slot
func (s *expertAvailabilityService) UpdateAvailability(id string, req *model.UpdateAvailabilityRequest) (*model.Availability, error) {
	// TODO: Implement update in Redis
	// This might need to fetch all availabilities for the expert and date,
	// update the specific one, and save back
	return nil, fmt.Errorf("update availability not implemented with Redis")
}

// DeleteAvailability deletes an availability slot
func (s *expertAvailabilityService) DeleteAvailability(id string) error {
	// TODO: Implement delete in Redis
	return fmt.Errorf("delete availability not implemented with Redis")
}

// GetAvailabilities retrieves filtered availability slots
func (s *expertAvailabilityService) GetAvailabilities(expertID string, startDate, endDate time.Time, isBooked *bool) ([]*model.Availability, error) {
	var availabilities []*model.Availability

	// Lặp qua từng ngày trong khoảng thời gian
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		key := fmt.Sprintf("availability:%s:%s", expertID, d.Format("2006-01-02"))
		data, err := s.cache.GetAvailability(key)
		if err != nil || data == nil {
			continue // Skip if error or not found, try next date
		}

		var availability model.Availability
		if err := json.Unmarshal(data, &availability); err != nil {
			continue // Skip if error, try next date
		}

		if isBooked == nil || availability.IsBooked == *isBooked {
			availabilities = append(availabilities, &availability)
		}
	}

	return availabilities, nil
}

// BookAvailability books an availability slot
func (s *expertAvailabilityService) BookAvailability(id string) error {
	// TODO: Implement booking in Redis
	// This might need to fetch all availabilities for the expert and date,
	// update the specific one's IsBooked status, and save back
	return fmt.Errorf("book availability not implemented with Redis")
}

// CreateRecurringAvailability creates multiple availability slots for recurring schedules
func (s *expertAvailabilityService) CreateRecurringAvailability(req *model.CreateRecurringAvailabilityRequest) ([]*model.Availability, error) {
	expertUUID, err := uuid.Parse(req.ExpertID)
	if err != nil {
		return nil, fmt.Errorf("invalid expert ID format: %v", err)
	}
	expert, err := s.expertRepo.GetByID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
	}
	if expert == nil {
		return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %s", req.ExpertID)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("định dạng ngày bắt đầu không hợp lệ")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("định dạng ngày kết thúc không hợp lệ")
	}

	if endDate.Before(startDate) {
		return nil, fmt.Errorf("ngày kết thúc phải sau ngày bắt đầu")
	}

	var createdAvailabilities []*model.Availability
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		weekday := int(d.Weekday())
		for _, w := range req.DaysOfWeek {
			if weekday == w {
				availability := &model.Availability{
					ID:        uuid.New(),
					ExpertID:  req.ExpertID,
					Date:      d.Format("2006-01-02"),
					StartTime: req.StartTime,
					EndTime:   req.EndTime,
					IsBooked:  false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				// Lưu vào Redis
				key := fmt.Sprintf("availability:%s:%s", req.ExpertID, d.Format("2006-01-02"))
				data, err := json.Marshal(availability)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal availability: %w", err)
				}

				if err := s.cache.SetAvailability(key, data); err != nil {
					return nil, fmt.Errorf("failed to cache availability: %w", err)
				}

				createdAvailabilities = append(createdAvailabilities, availability)
			}
		}
	}
	return createdAvailabilities, nil
}

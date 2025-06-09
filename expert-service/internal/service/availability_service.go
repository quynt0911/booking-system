package service

import (
    "expert-service/internal/cache"
    "expert-service/internal/model"
    "expert-service/internal/repository"
    "fmt"
    "time"
)

type ExpertAvailabilityService interface {
    CheckAvailability(req *model.CheckAvailabilityRequest) (bool, error)
    CreateOffTime(req *model.CreateOffTimeRequest) (*model.OffTime, error)
    GetExpertOffTimes(expertID int) ([]*model.OffTime, error)
    DeleteOffTime(id int) error
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
    requestDate, err := time.Parse("2006-01-02", req.Date)
    if err != nil {
        return false, fmt.Errorf("định dạng ngày không hợp lệ (sử dụng YYYY-MM-DD)")
    }

    requestTime, err := time.Parse("15:04", req.Time)
    if err != nil {
        return false, fmt.Errorf("định dạng thời gian không hợp lệ (sử dụng HH:MM)")
    }

    cacheKey := req.Date
    if isAvailable, exists, err := s.cache.GetAvailability(req.ExpertID, cacheKey); err == nil && exists {
        return isAvailable, nil
    }

    expert, err := s.expertRepo.GetByID(req.ExpertID)
    if err != nil {
        return false, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
    }
    if expert == nil {
        return false, fmt.Errorf("không tìm thấy chuyên gia với ID %d", req.ExpertID)
    }
    if expert.Status != "active" {
        return false, nil
    }

    offTimes, err := s.offTimeRepo.GetByExpertIDAndDateRange(req.ExpertID, requestDate)
    if err != nil {
        return false, fmt.Errorf("không thể kiểm tra thời gian nghỉ: %v", err)
    }
    if len(offTimes) > 0 {
        _ = s.cache.SetAvailability(req.ExpertID, cacheKey, false)
        return false, nil
    }

    dayOfWeek := int(requestDate.Weekday())
    schedules, err := s.scheduleRepo.GetByExpertIDAndDay(req.ExpertID, dayOfWeek)
    if err != nil {
        return false, fmt.Errorf("không thể kiểm tra lịch trình: %v", err)
    }

    isAvailable := false
    for _, schedule := range schedules {
        startTime, err1 := time.Parse("15:04", schedule.StartTime)
        endTime, err2 := time.Parse("15:04", schedule.EndTime)
        if err1 != nil || err2 != nil {
            continue // Bỏ qua lịch trình lỗi định dạng
        }
        if (requestTime.Equal(startTime) || requestTime.After(startTime)) && requestTime.Before(endTime) {
            isAvailable = true
            break
        }
    }

    _ = s.cache.SetAvailability(req.ExpertID, cacheKey, isAvailable)
    return isAvailable, nil
}

// CreateOffTime tạo thời gian nghỉ cho chuyên gia
func (s *expertAvailabilityService) CreateOffTime(req *model.CreateOffTimeRequest) (*model.OffTime, error) {
    expert, err := s.expertRepo.GetByID(req.ExpertID)
    if err != nil {
        return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
    }
    if expert == nil {
        return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %d", req.ExpertID)
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

    offTime := &model.OffTime{
        ExpertID:  req.ExpertID,
        StartDate: startDate,
        EndDate:   endDate,
        Reason:    req.Reason,
    }

	err = s.offTimeRepo.Create(offTime)
    if err != nil {
        return nil, fmt.Errorf("không thể tạo thời gian nghỉ: %v", err)
    }

    _ = s.cache.InvalidateExpert(req.ExpertID)
    return offTime, nil
}

// GetExpertOffTimes lấy danh sách thời gian nghỉ của chuyên gia
func (s *expertAvailabilityService) GetExpertOffTimes(expertID int) ([]*model.OffTime, error) {
    expert, err := s.expertRepo.GetByID(expertID)
    if err != nil {
        return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
    }
    if expert == nil {
        return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %d", expertID)
    }

    offTimes, err := s.offTimeRepo.GetByExpertID(expertID)
    if err != nil {
        return nil, fmt.Errorf("không thể lấy danh sách thời gian nghỉ: %v", err)
    }

    return offTimes, nil
}

// DeleteOffTime xóa thời gian nghỉ
func (s *expertAvailabilityService) DeleteOffTime(id int) error {
    err := s.offTimeRepo.Delete(id)
    if err != nil {
        return fmt.Errorf("không thể xóa thời gian nghỉ: %v", err)
    }
    return nil
}
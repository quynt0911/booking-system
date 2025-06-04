package service

import (
    "expert-service/internal2/cache"
    "expert-service/internal2/model"
    "expert-service/internal2/repository"
    "fmt"
    "time"
)

type AvailabilityService interface {
    CheckAvailability(req *model.CheckAvailabilityRequest) (bool, error)
    CreateOffTime(req *model.CreateOffTimeRequest) (*model.OffTime, error)
    GetExpertOffTimes(expertID int) ([]*model.OffTime, error)
    DeleteOffTime(id int) error
}

type availabilityService struct {
    expertRepo   repository.ExpertRepository
    scheduleRepo repository.ScheduleRepository
    offTimeRepo  repository.OffTimeRepository
    cache        cache.AvailabilityCache
}

func NewAvailabilityService(
    expertRepo repository.ExpertRepository,
    scheduleRepo repository.ScheduleRepository,
    offTimeRepo repository.OffTimeRepository,
    cache cache.AvailabilityCache,
) AvailabilityService {
    return &availabilityService{
        expertRepo:   expertRepo,
        scheduleRepo: scheduleRepo,
        offTimeRepo:  offTimeRepo,
        cache:        cache,
    }
}

func (s *availabilityService) CheckAvailability(req *model.CheckAvailabilityRequest) (bool, error) {
    // Parse request date
    requestDate, err := time.Parse("2006-01-02", req.Date)
    if err != nil {
        return false, fmt.Errorf("định dạng ngày không hợp lệ (sử dụng YYYY-MM-DD)")
    }
    
    // Parse request time
    requestTime, err := time.Parse("15:04", req.Time)
    if err != nil {
        return false, fmt.Errorf("định dạng thời gian không hợp lệ (sử dụng HH:MM)")
    }
    
    // Check cache first
    cacheKey := req.Date
    if isAvailable, exists, err := s.cache.GetAvailability(req.ExpertID, cacheKey); err == nil && exists {
        return isAvailable, nil
    }
    
    // 1. Check if expert exists and is active
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
    
    // 2. Check if expert has off time on this date
    offTimes, err := s.offTimeRepo.GetByExpertIDAndDateRange(req.ExpertID, requestDate)
    if err != nil {
        return false, fmt.Errorf("không thể kiểm tra thời gian nghỉ: %v", err)
    }
    if len(offTimes) > 0 {
        // Expert is on leave
        s.cache.SetAvailability(req.ExpertID, cacheKey, false)
        return false, nil
    }
    
    // 3. Check if expert has schedule on this day of week
    dayOfWeek := int(requestDate.Weekday())
    schedules, err := s.scheduleRepo.GetByExpertIDAndDay(req.ExpertID, dayOfWeek)
    if err != nil {
        return false, fmt.Errorf("không thể kiểm tra lịch trình: %v", err)
    }
    
    // 4. Check if request time falls within any schedule
    isAvailable := false
    for _, schedule := range schedules {
        startTime, _ := time.Parse("15:04", schedule.StartTime)
        endTime, _ := time.Parse("15:04", schedule.EndTime)
        
        if (requestTime.Equal(startTime) || requestTime.After(startTime)) && requestTime.Before(endTime) {
            isAvailable = true
            break
        }
    }
    
    // Cache the result
    s.cache.SetAvailability(req.ExpertID, cacheKey, isAvailable)
    
    return isAvailable, nil
}

func (s *availabilityService) CreateOffTime(req *model.CreateOffTimeRequest) (*model.OffTime, error) {
    // Validate expert exists
    expert, err := s.expertRepo.GetByID(req.ExpertID)
    if err != nil {
        return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
    }
    if expert == nil {
        return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %d", req.ExpertID)
    }
    
    // Parse dates
    startDate, err := time.Parse("2006-01-02", req.StartDate)
    if err != nil {
        return nil, fmt.Errorf("định dạng ngày bắt đầu không hợp lệ")
    }
    
    endDate, err := time.Parse("2006-01-02", req.EndDate)
    if err != nil {
        return nil, fmt.Errorf("định dạng ngày kết thúc không hợp lệ")
    }
    
    // Validate date range
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
    
    // Invalidate cache for affected dates
    s.cache.InvalidateExpert(req.ExpertID)
    
    return offTime, nil
}

func (s *availabilityService) GetExpertOffTimes(expertID int) ([]*model.OffTime, error) {
    // Validate expert exists
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

func (s *availabilityService) DeleteOffTime(id int) error {
    err := s.offTimeRepo.Delete(id)
    if err != nil {
        return fmt.Errorf("không thể xóa thời gian nghỉ: %v", err)
    }
    
    return nil
}
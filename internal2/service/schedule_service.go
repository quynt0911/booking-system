package service

import (
    "expert-service/internal2/cache"
    "expert-service/internal2/model"
    "expert-service/internal2/repository"
    "fmt"
    "time"
)

type ScheduleService interface {
    CreateSchedule(req *model.CreateScheduleRequest) (*model.Schedule, error)
    GetExpertSchedules(expertID int) ([]*model.Schedule, error)
    UpdateSchedule(id int, req *model.CreateScheduleRequest) (*model.Schedule, error)
    DeleteSchedule(id int) error
}

type scheduleService struct {
    scheduleRepo repository.ScheduleRepository
    expertRepo   repository.ExpertRepository
    cache        cache.AvailabilityCache
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, expertRepo repository.ExpertRepository, cache cache.AvailabilityCache) ScheduleService {
    return &scheduleService{
        scheduleRepo: scheduleRepo,
        expertRepo:   expertRepo,
        cache:        cache,
    }
}

func (s *scheduleService) CreateSchedule(req *model.CreateScheduleRequest) (*model.Schedule, error) {
    // Validate expert exists
    expert, err := s.expertRepo.GetByID(req.ExpertID)
    if err != nil {
        return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
    }
    if expert == nil {
        return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %d", req.ExpertID)
    }
    
    // Validate time format
    if !s.isValidTimeFormat(req.StartTime) || !s.isValidTimeFormat(req.EndTime) {
        return nil, fmt.Errorf("định dạng thời gian không hợp lệ (sử dụng HH:MM)")
    }
    
    // Validate time logic
    if !s.isValidTimeRange(req.StartTime, req.EndTime) {
        return nil, fmt.Errorf("thời gian bắt đầu phải nhỏ hơn thời gian kết thúc")
    }
    
    schedule := &model.Schedule{
        ExpertID:  req.ExpertID,
        DayOfWeek: req.DayOfWeek,
        StartTime: req.StartTime,
        EndTime:   req.EndTime,
        IsActive:  true,
    }
    
    err = s.scheduleRepo.Create(schedule)
    if err != nil {
        return nil, fmt.Errorf("không thể tạo lịch trình: %v", err)
    }
    
    // Invalidate cache when schedule changes
    s.cache.InvalidateExpert(req.ExpertID)
    
    return schedule, nil
}

func (s *scheduleService) GetExpertSchedules(expertID int) ([]*model.Schedule, error) {
    // Validate expert exists
    expert, err := s.expertRepo.GetByID(expertID)
    if err != nil {
        return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
    }
    if expert == nil {
        return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %d", expertID)
    }
    
    schedules, err := s.scheduleRepo.GetByExpertID(expertID)
    if err != nil {
        return nil, fmt.Errorf("không thể lấy lịch trình: %v", err)
    }
    
    return schedules, nil
}

func (s *scheduleService) UpdateSchedule(id int, req *model.CreateScheduleRequest) (*model.Schedule, error) {
    // Validate time format and logic
    if !s.isValidTimeFormat(req.StartTime) || !s.isValidTimeFormat(req.EndTime) {
        return nil, fmt.Errorf("định dạng thời gian không hợp lệ")
    }
    
    if !s.isValidTimeRange(req.StartTime, req.EndTime) {
        return nil, fmt.Errorf("thời gian bắt đầu phải nhỏ hơn thời gian kết thúc")
    }
    
    schedule := &model.Schedule{
        ID:        id,
        ExpertID:  req.ExpertID,
        DayOfWeek: req.DayOfWeek,
        StartTime: req.StartTime,
        EndTime:   req.EndTime,
        IsActive:  true,
    }
    
    err := s.scheduleRepo.Update(schedule)
    if err != nil {
        return nil, fmt.Errorf("không thể cập nhật lịch trình: %v", err)
    }
    
    // Invalidate cache
    s.cache.InvalidateExpert(req.ExpertID)
    
    return schedule, nil
}

func (s *scheduleService) DeleteSchedule(id int) error {
    err := s.scheduleRepo.Delete(id)
    if err != nil {
        return fmt.Errorf("không thể xóa lịch trình: %v", err)
    }
    
    return nil
}

// Helper functions
func (s *scheduleService) isValidTimeFormat(timeStr string) bool {
    _, err := time.Parse("15:04", timeStr)
    return err == nil
}

func (s *scheduleService) isValidTimeRange(startTime, endTime string) bool {
    start, err1 := time.Parse("15:04", startTime)
    end, err2 := time.Parse("15:04", endTime)
    
    if err1 != nil || err2 != nil {
        return false
    }
    
    return start.Before(end)
}
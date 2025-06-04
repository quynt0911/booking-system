package service

import (
    "expert-service/internal2/cache"
    "expert-service/internal2/model"
    "expert-service/internal2/repository"
    "fmt"
)

type ExpertService interface {
    CreateExpert(req *model.CreateExpertRequest) (*model.Expert, error)
    GetExpert(id int) (*model.Expert, error)
    GetAllExperts() ([]*model.Expert, error)
    UpdateExpert(id int, req *model.CreateExpertRequest) (*model.Expert, error)
    DeleteExpert(id int) error
}

type expertService struct {
    expertRepo repository.ExpertRepository
    cache      cache.AvailabilityCache
}

func NewExpertService(expertRepo repository.ExpertRepository, cache cache.AvailabilityCache) ExpertService {
    return &expertService{
        expertRepo: expertRepo,
        cache:      cache,
    }
}

func (s *expertService) CreateExpert(req *model.CreateExpertRequest) (*model.Expert, error) {
    // Validate business rules
    if req.Name == "" {
        return nil, fmt.Errorf("tên chuyên gia không được để trống")
    }
    
    expert := &model.Expert{
        Name:           req.Name,
        Email:          req.Email,
        Specialization: req.Specialization,
        Status:         "active",
    }
    
    err := s.expertRepo.Create(expert)
    if err != nil {
        return nil, fmt.Errorf("không thể tạo chuyên gia: %v", err)
    }
    
    return expert, nil
}

func (s *expertService) GetExpert(id int) (*model.Expert, error) {
    expert, err := s.expertRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("không thể lấy thông tin chuyên gia: %v", err)
    }
    
    if expert == nil {
        return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %d", id)
    }
    
    return expert, nil
}

func (s *expertService) GetAllExperts() ([]*model.Expert, error) {
    experts, err := s.expertRepo.GetAll()
    if err != nil {
        return nil, fmt.Errorf("không thể lấy danh sách chuyên gia: %v", err)
    }
    
    return experts, nil
}

func (s *expertService) UpdateExpert(id int, req *model.CreateExpertRequest) (*model.Expert, error) {
    // Check if expert exists
    existingExpert, err := s.expertRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("không thể tìm chuyên gia: %v", err)
    }
    
    if existingExpert == nil {
        return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %d", id)
    }
    
    // Update fields
    existingExpert.Name = req.Name
    existingExpert.Email = req.Email
    existingExpert.Specialization = req.Specialization
    
    err = s.expertRepo.Update(existingExpert)
    if err != nil {
        return nil, fmt.Errorf("không thể cập nhật chuyên gia: %v", err)
    }
    
    // Invalidate cache when expert info changes
    s.cache.InvalidateExpert(id)
    
    return existingExpert, nil
}

func (s *expertService) DeleteExpert(id int) error {
    // Check if expert exists
    expert, err := s.expertRepo.GetByID(id)
    if err != nil {
        return fmt.Errorf("không thể tìm chuyên gia: %v", err)
    }
    
    if expert == nil {
        return fmt.Errorf("không tìm thấy chuyên gia với ID %d", id)
    }
    
    err = s.expertRepo.Delete(id)
    if err != nil {
        return fmt.Errorf("không thể xóa chuyên gia: %v", err)
    }
    
    // Invalidate cache when expert is deleted
    s.cache.InvalidateExpert(id)
    
    return nil
}
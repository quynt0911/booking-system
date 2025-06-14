package service

import (
	"expert-service/internal/model"
	"expert-service/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ExpertService interface {
	CreateExpert(req *model.CreateExpertRequest) (*model.Expert, error)
	GetExpertByID(id uuid.UUID) (*model.Expert, error)
	GetExpertByEmail(email string) (*model.Expert, error)
	GetAllExperts(limit, offset int) ([]*model.Expert, error)
	UpdateExpert(id uuid.UUID, req *model.UpdateExpertRequest) error
	DeleteExpert(id uuid.UUID) error
	GetExpertsByExpertise(expertise string) ([]*model.Expert, error)
}

type expertService struct {
	expertRepo repository.ExpertRepository
}

func NewExpertService(expertRepo repository.ExpertRepository) ExpertService {
	return &expertService{
		expertRepo: expertRepo,
	}
}

func (s *expertService) CreateExpert(req *model.CreateExpertRequest) (*model.Expert, error) {
	// Check if expert with user_id already exists
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	expert := &model.Expert{
		ID:              uuid.New(),
		UserID:          userID,
		Specialization:  req.Specialization,
		ExperienceYears: req.ExperienceYears,
		HourlyRate:      req.HourlyRate,
		Certifications:  req.Certifications,
		IsAvailable:     req.IsAvailable,
		Rating:          0,
		TotalReviews:    0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.expertRepo.Create(expert); err != nil {
		return nil, fmt.Errorf("failed to create expert: %w", err)
	}

	return expert, nil
}

func (s *expertService) GetExpertByID(id uuid.UUID) (*model.Expert, error) {
	expert, err := s.expertRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get expert: %w", err)
	}
	if expert == nil {
		return nil, fmt.Errorf("expert not found")
	}
	return expert, nil
}

func (s *expertService) GetExpertByEmail(email string) (*model.Expert, error) {
	return s.expertRepo.GetByEmail(email)
}

func (s *expertService) GetAllExperts(limit, offset int) ([]*model.Expert, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.expertRepo.GetAll()
}

func (s *expertService) UpdateExpert(id uuid.UUID, req *model.UpdateExpertRequest) error {
	// Check if expert exists
	expert, err := s.expertRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get expert: %w", err)
	}
	if expert == nil {
		return fmt.Errorf("expert not found")
	}

	// Update expert fields
	if req.Specialization != nil {
		expert.Specialization = *req.Specialization
	}
	if req.ExperienceYears != nil {
		expert.ExperienceYears = *req.ExperienceYears
	}
	if req.HourlyRate != nil {
		expert.HourlyRate = *req.HourlyRate
	}
	if req.Certifications != nil {
		expert.Certifications = req.Certifications
	}
	if req.IsAvailable != nil {
		expert.IsAvailable = *req.IsAvailable
	}
	expert.UpdatedAt = time.Now()

	return s.expertRepo.Update(expert)
}

func (s *expertService) DeleteExpert(id uuid.UUID) error {
	// Check if expert exists
	expert, err := s.expertRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get expert: %w", err)
	}
	if expert == nil {
		return fmt.Errorf("expert not found")
	}

	return s.expertRepo.Delete(id)
}

func (s *expertService) GetExpertsByExpertise(expertise string) ([]*model.Expert, error) {
	return s.expertRepo.GetByExpertise(expertise)
}

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
	// Check if expert with email already exists
	existing, err := s.expertRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing expert: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("expert with email %s already exists", req.Email)
	}

	expert := &model.Expert{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Expertise: req.Specialization,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

	// Check email uniqueness if email is being updated
	if req.Email != nil && *req.Email != expert.Email {
		existing, err := s.expertRepo.GetByEmail(*req.Email)
		if err != nil {
			return fmt.Errorf("failed to check existing expert: %w", err)
		}
		if existing != nil {
			return fmt.Errorf("expert with email %s already exists", *req.Email)
		}
	}

	return s.expertRepo.Update(id, req)
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
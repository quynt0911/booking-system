package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"services/booking-service/internal/model"
	"services/booking-service/internal/repository"
	"services/booking-service/pkg/logger"
)

type StatusServiceInterface interface {
	UpdateBookingStatus(bookingID uuid.UUID, status model.BookingStatus, changedBy uuid.UUID, userRole, note string) error
	GetBookingStatus(bookingID uuid.UUID) (model.BookingStatus, error)
	GetStatusHistory(bookingID uuid.UUID, userID uuid.UUID, userRole string) ([]*model.StatusHistory, error)
	ValidateStatusTransition(currentStatus, newStatus model.BookingStatus, userRole string) error
}

type StatusService struct {
	statusHistoryRepo repository.StatusHistoryRepositoryInterface
	bookingRepo       repository.BookingRepositoryInterface
	logger            logger.LoggerInterface
}

func NewStatusService(
	statusHistoryRepo repository.StatusHistoryRepositoryInterface,
	bookingRepo repository.BookingRepositoryInterface,
	logger logger.LoggerInterface,
) StatusServiceInterface {
	return &StatusService{
		statusHistoryRepo: statusHistoryRepo,
		bookingRepo:       bookingRepo,
		logger:            logger,
	}
}

func (s *StatusService) UpdateBookingStatus(bookingID uuid.UUID, status model.BookingStatus, changedBy uuid.UUID, userRole, note string) error {
	// Get current booking
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return fmt.Errorf("booking not found")
	}

	// Check authorization
	if err := s.checkStatusUpdateAuthorization(booking, changedBy, userRole); err != nil {
		return err
	}

	// Validate status transition
	if err := s.ValidateStatusTransition(booking.Status, status, userRole); err != nil {
		return err
	}

	// Update booking status
	booking.Status = status
	booking.UpdatedAt = time.Now()

	_, err = s.bookingRepo.Update(booking)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %v", err)
	}

	// Create status history record
	statusHistory := &model.StatusHistory{
		BookingID: bookingID,
		Status:    status,
		ChangedBy: changedBy,
		ChangedAt: time.Now(),
		Note:      note,
	}

	if err := s.statusHistoryRepo.Create(statusHistory); err != nil {
		s.logger.Error("Failed to create status history", err)
		// Don't return error as the main update succeeded
	}

	s.logger.Info(fmt.Sprintf("Booking %s status updated to %s by user %s", bookingID, status, changedBy))

	return nil
}

func (s *StatusService) GetBookingStatus(bookingID uuid.UUID) (model.BookingStatus, error) {
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return "", fmt.Errorf("booking not found")
	}

	return booking.Status, nil
}

func (s *StatusService) GetStatusHistory(bookingID uuid.UUID, userID uuid.UUID, userRole string) ([]*model.StatusHistory, error) {
	// Get booking to check authorization
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return nil, fmt.Errorf("booking not found")
	}

	// Check authorization
	if userRole != "admin" &&
		booking.UserID != userID &&
		booking.ExpertID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Get status history
	history, err := s.statusHistoryRepo.GetByBookingID(bookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get status history: %v", err)
	}

	// Convert []model.StatusHistory to []*model.StatusHistory
	historyPtrs := make([]*model.StatusHistory, len(history))
	for i := range history {
		historyPtrs[i] = &history[i]
	}
	return historyPtrs, nil
}

func (s *StatusService) ValidateStatusTransition(currentStatus, newStatus model.BookingStatus, userRole string) error {
	// Define valid transitions
	validTransitions := map[model.BookingStatus][]model.BookingStatus{
		model.BookingStatusPending:   {model.BookingStatusConfirmed, model.BookingStatusCancelled},
		model.BookingStatusConfirmed: {model.BookingStatusCompleted, model.BookingStatusCancelled},
		model.BookingStatusCompleted: {},
		model.BookingStatusCancelled: {},
	}

	// Admin can change to any status
	if userRole == "admin" {
		return nil
	}

	// Check if the transition is valid
	if allowedStatuses, exists := validTransitions[currentStatus]; exists {
		for _, status := range allowedStatuses {
			if status == newStatus {
				return nil
			}
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
}

// checkStatusUpdateAuthorization checks if the user is authorized to update the booking status.
func (s *StatusService) checkStatusUpdateAuthorization(booking *model.Booking, userID uuid.UUID, userRole string) error {
	if userRole == "admin" {
		return nil
	}
	if booking.UserID == userID || booking.ExpertID == userID {
		return nil
	}
	return fmt.Errorf("user is not authorized to update this booking status")
}

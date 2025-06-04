package service

import (
	"fmt"
	"time"

	"booking-service/internal/model"
	"booking-service/internal/repository"
	"shared/pkg/logger"
)

type StatusServiceInterface interface {
	UpdateBookingStatus(bookingID uint, status string, changedBy uint, userRole, note string) error
	GetBookingStatus(bookingID uint) (string, error)
	GetStatusHistory(bookingID uint, userID uint, userRole string) ([]*model.StatusHistory, error)
	ValidateStatusTransition(currentStatus, newStatus, userRole string) error
}

type StatusService struct {
	statusHistoryRepo repository.StatusHistoryRepositoryInterface
	bookingRepo       repository.BookingRepositoryInterface
	logger            logger.Logger
}

func NewStatusService(
	statusHistoryRepo repository.StatusHistoryRepositoryInterface,
	logger logger.Logger,
) StatusServiceInterface {
	return &StatusService{
		statusHistoryRepo: statusHistoryRepo,
		logger:            logger,
	}
}

func (s *StatusService) UpdateBookingStatus(bookingID uint, status string, changedBy uint, userRole, note string) error {
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

	s.logger.Info(fmt.Sprintf("Booking %d status updated to %s by user %d", bookingID, status, changedBy))

	return nil
}

func (s *StatusService) GetBookingStatus(bookingID uint) (string, error) {
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return "", fmt.Errorf("booking not found")
	}

	return booking.Status, nil
}

func (s *StatusService) GetStatusHistory(bookingID uint, userID uint, userRole string) ([]*model.StatusHistory, error) {
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

	return history, nil
}

func (s *StatusService) ValidateStatusTransition(currentStatus, newStatus, userRole string) error {
	// Define valid transitions
	validTransitions := map[string][]string{
		model.BookingStatusPending:    {model.BookingStatusConfirmed, model.BookingStatusCancelled},
		model.BookingStatusConfirmed:  {model.BookingStatusCompleted, model.BookingStatusCancelled},
		model.BookingStatusCompleted:  {},
		model.BookingStatusCancelled:  {},
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
func (s *StatusService) checkStatusUpdateAuthorization(booking *model.Booking, userID uint, userRole string) error {
	if userRole == "admin" {
		return nil
	}
	if booking.UserID == userID || booking.ExpertID == userID {
		return nil
	}
	return fmt.Errorf("user is not authorized to update this booking status")
}
package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"booking-service/internal/model"
	"booking-service/internal/repository"
	"shared/pkg/logger"
)

type BookingServiceInterface interface {
	CreateBooking(userID uint, req *model.CreateBookingRequest) (*model.Booking, error)
	GetBookingByID(bookingID uint) (*model.Booking, error)
	UpdateBooking(bookingID uint, req *model.UpdateBookingRequest) (*model.Booking, error)
	CancelBooking(bookingID uint, userID uint) error
	GetUserBookings(userID uint, page, limit int, status string) ([]*model.Booking, int, error)
	GetExpertBookings(expertID uint, page, limit int, status, date string) ([]*model.Booking, int, error)
	GetBookingHistory(userID uint, page, limit int, status string, startDate, endDate *time.Time) ([]*model.Booking, int, error)
	GetExpertHistory(expertID uint, page, limit int, status string, startDate, endDate *time.Time) ([]*model.Booking, int, error)
	GetBookingStatistics(userID uint, userRole, period string, year, month int) (map[string]interface{}, error)
	GetUpcomingUserBookings(userID uint, limit, days int) ([]*model.Booking, error)
	GetUpcomingExpertBookings(expertID uint, limit, days int) ([]*model.Booking, error)
	GetPastUserBookings(userID uint, page, limit int) ([]*model.Booking, int, error)
	GetPastExpertBookings(expertID uint, page, limit int) ([]*model.Booking, int, error)
}

type BookingService struct {
	bookingRepo       repository.BookingRepositoryInterface
	statusHistoryRepo repository.StatusHistoryRepositoryInterface
	redisClient       *redis.Client
	logger            logger.Logger
}

func NewBookingService(
	bookingRepo repository.BookingRepositoryInterface,
	statusHistoryRepo repository.StatusHistoryRepositoryInterface,
	redisClient *redis.Client,
	logger logger.Logger,
) BookingServiceInterface {
	return &BookingService{
		bookingRepo:       bookingRepo,
		statusHistoryRepo: statusHistoryRepo,
		redisClient:       redisClient,
		logger:            logger,
	}
}

func (s *BookingService) CreateBooking(userID uint, req *model.CreateBookingRequest) (*model.Booking, error) {
	// Create booking model
	booking := &model.Booking{
		UserID:      userID,
		ExpertID:    req.ExpertID,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Type:        req.Type,
		Note:        req.Note,
		Status:      model.BookingStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save booking to database
	createdBooking, err := s.bookingRepo.Create(booking)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %v", err)
	}

	// Create initial status history
	statusHistory := &model.StatusHistory{
		BookingID:   createdBooking.ID,
		Status:      model.BookingStatusPending,
		ChangedBy:   userID,
		ChangedAt:   time.Now(),
		Note:        "Booking created",
	}

	if err := s.statusHistoryRepo.Create(statusHistory); err != nil {
		s.logger.Error("Failed to create status history", err)
	}

	// Cache booking data
	s.cacheBooking(createdBooking)

	// TODO: Send notification to expert
	s.notifyBookingCreated(createdBooking)

	return createdBooking, nil
}

func (s *BookingService) GetBookingByID(bookingID uint) (*model.Booking, error) {
	// Try to get from cache first
	if booking := s.getBookingFromCache(bookingID); booking != nil {
		return booking, nil
	}

	// Get from database
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.cacheBooking(booking)

	return booking, nil
}

func (s *BookingService) UpdateBooking(bookingID uint, req *model.UpdateBookingRequest) (*model.Booking, error) {
	// Get existing booking
	booking, err := s.GetBookingByID(bookingID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.StartTime != nil {
		booking.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		booking.EndTime = *req.EndTime
	}
	if req.Note != nil {
		booking.Note = *req.Note
	}
	if req.Type != nil {
		booking.Type = *req.Type
	}

	booking.UpdatedAt = time.Now()

	// Save to database
	updatedBooking, err := s.bookingRepo.Update(booking)
	if err != nil {
		return nil, fmt.Errorf("failed to update booking: %v", err)
	}

	// Update cache
	s.cacheBooking(updatedBooking)

	// TODO: Send notification about update
	s.notifyBookingUpdated(updatedBooking)

	return updatedBooking, nil
}

func (s *BookingService) CancelBooking(bookingID uint, userID uint) error {
	// Get existing booking
	booking, err := s.GetBookingByID(bookingID)
	if err != nil {
		return err
	}

	// Update status to cancelled
	booking.Status = model.BookingStatusCancelled
	booking.UpdatedAt = time.Now()

	// Save to database
	_, err = s.bookingRepo.Update(booking)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %v", err)
	}

	// Create status history
	statusHistory := &model.StatusHistory{
		BookingID:   bookingID,
		Status:      model.BookingStatusCancelled,
		ChangedBy:   userID,
		ChangedAt:   time.Now(),
		Note:        "Booking cancelled by user",
	}

	if err := s.statusHistoryRepo.Create(statusHistory); err != nil {
		s.logger.Error("Failed to create status history", err)
	}

	// Update cache
	s.cacheBooking(booking)

	// TODO: Send notification about cancellation
	s.notifyBookingCancelled(booking)

	return nil
}

func (s *BookingService) GetUserBookings(userID uint, page, limit int, status string) ([]*model.Booking, int, error) {
	offset := (page - 1) * limit
	return s.bookingRepo.GetByUserID(userID, offset, limit, status)
}

func (s *BookingService) GetExpertBookings(expertID uint, page, limit int, status, date string) ([]*model.Booking, int, error) {
	offset := (page - 1) * limit
	return s.bookingRepo.GetByExpertID(expertID, offset, limit, status, date)
}

func (s *BookingService) GetBookingHistory(userID uint, page, limit int, status string, startDate, endDate *time.Time) ([]*model.Booking, int, error) {
	offset := (page - 1) * limit
	return s.bookingRepo.GetHistoryByUserID(userID, offset, limit, status, startDate, endDate)
}

func (s *BookingService) GetExpertHistory(expertID uint, page, limit int, status string, startDate, endDate *time.Time) ([]*model.Booking, int, error) {
	offset := (page - 1) * limit
	return s.bookingRepo.GetHistoryByExpertID(expertID, offset, limit, status, startDate, endDate)
}

func (s *BookingService) GetBookingStatistics(userID uint, userRole, period string, year, month int) (map[string]interface{}, error) {
	var stats map[string]interface{}
	var err error

	if userRole == "expert" {
		stats, err = s.bookingRepo.GetExpertStatistics(userID, period, year, month)
	} else {
		stats, err = s.bookingRepo.GetUserStatistics(userID, period, year, month)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %v", err)
	}

	return stats, nil
}

func (s *BookingService) GetUpcomingUserBookings(userID uint, limit, days int) ([]*model.Booking, error) {
	endDate := time.Now().AddDate(0, 0, days)
	return s.bookingRepo.GetUpcomingByUserID(userID, limit, endDate)
}

func (s *BookingService) GetUpcomingExpertBookings(expertID uint, limit, days int) ([]*model.Booking, error) {
	endDate := time.Now().AddDate(0, 0, days)
	return s.bookingRepo.GetUpcomingByExpertID(expertID, limit, endDate)
}

func (s *BookingService) GetPastUserBookings(userID uint, page, limit int) ([]*model.Booking, int, error) {
	offset := (page - 1) * limit
	return s.bookingRepo.GetPastByUserID(userID, offset, limit)
}

func (s *BookingService) GetPastExpertBookings(expertID uint, page, limit int) ([]*model.Booking, int, error) {
	offset := (page - 1) * limit
	return s.bookingRepo.GetPastByExpertID(expertID, offset, limit)
}

// Cache helper methods
func (s *BookingService) cacheBooking(booking *model.Booking) {
	key := fmt.Sprintf("booking:%d", booking.ID)
	data, err := json.Marshal(booking)
	if err != nil {
		s.logger.Error("Failed to marshal booking for cache", err)
		return
	}

	err = s.redisClient.Set(s.redisClient.Context(), key, data, 15*time.Minute).Err()
	if err != nil {
		s.logger.Error("Failed to cache booking", err)
	}
}

func (s *BookingService) getBookingFromCache(bookingID uint) *model.Booking {
	key := fmt.Sprintf("booking:%d", bookingID)
	data, err := s.redisClient.Get(s.redisClient.Context(), key).Result()
	if err != nil {
		return nil
	}

	var booking model.Booking
	if err := json.Unmarshal([]byte(data), &booking); err != nil {
		s.logger.Error("Failed to unmarshal cached booking", err)
		return nil
	}

	return &booking
}

// Notification helper methods (placeholders)
func (s *BookingService) notifyBookingCreated(booking *model.Booking) {
	// TODO: Integrate with notification service
	s.logger.Info(fmt.Sprintf("Booking created: %d", booking.ID))
}

func (s *BookingService) notifyBookingUpdated(booking *model.Booking) {
	// TODO: Integrate with notification service
	s.logger.Info(fmt.Sprintf("Booking updated: %d", booking.ID))
}

func (s *BookingService) notifyBookingCancelled(booking *model.Booking) {
	// TODO: Integrate with notification service
	s.logger.Info(fmt.Sprintf("Booking cancelled: %d", booking.ID))
}
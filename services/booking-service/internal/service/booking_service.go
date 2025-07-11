package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"services/booking-service/internal/model"
	"services/booking-service/internal/repository"
	"services/booking-service/pkg/logger"
)

type BookingServiceInterface interface {
	CreateBooking(userID uuid.UUID, req *model.CreateBookingRequest) (*model.BookingResponse, error)
	GetBookingByID(bookingID uuid.UUID) (*model.BookingResponse, error)
	UpdateBooking(bookingID uuid.UUID, req *model.UpdateBookingRequest) (*model.BookingResponse, error)
	CancelBooking(bookingID uuid.UUID, userID uuid.UUID, req *model.CancelBookingRequest) error
	GetUserBookings(userID uuid.UUID, req *model.GetBookingsRequest) ([]model.BookingResponse, int64, error)
	GetExpertBookings(expertID uuid.UUID, req *model.GetBookingsRequest) ([]model.BookingResponse, int64, error)
	GetBookingHistory(userID uuid.UUID, req *model.GetHistoryRequest) ([]model.StatusHistoryResponse, int64, error)
	GetExpertHistory(expertID uuid.UUID, req *model.GetHistoryRequest) ([]model.StatusHistoryResponse, int64, error)
	GetBookingStatistics(userID uuid.UUID, userRole, period string, year, month int) (*model.BookingStatsResponse, error)
	GetUpcomingUserBookings(userID uuid.UUID, limit, days int) ([]model.BookingResponse, error)
	GetUpcomingExpertBookings(expertID uuid.UUID, limit, days int) ([]model.BookingResponse, error)
	GetPastUserBookings(userID uuid.UUID, page, limit int) ([]model.BookingResponse, int64, error)
	GetPastExpertBookings(expertID uuid.UUID, page, limit int) ([]model.BookingResponse, int64, error)
	CheckConflict(req *model.CheckConflictRequest) (*model.ConflictCheckResponse, error)
	CheckConflictWithExclusion(req *model.CheckConflictRequest, excludeID uuid.UUID) (*model.ConflictCheckResponse, error)
	GetExpertBookingsByDate(expertID uuid.UUID, date time.Time) ([]model.BookingResponse, error)
	GetExpertIDByUserID(userID uuid.UUID) (uuid.UUID, error)
	DeleteBooking(bookingID uuid.UUID) error
}

type BookingService struct {
	bookingRepo       repository.BookingRepositoryInterface
	statusHistoryRepo repository.StatusHistoryRepositoryInterface
	redisClient       *redis.Client
	logger            logger.LoggerInterface
}

func NewBookingService(
	bookingRepo repository.BookingRepositoryInterface,
	statusHistoryRepo repository.StatusHistoryRepositoryInterface,
	redisClient *redis.Client,
	logger logger.LoggerInterface,
) BookingServiceInterface {
	return &BookingService{
		bookingRepo:       bookingRepo,
		statusHistoryRepo: statusHistoryRepo,
		redisClient:       redisClient,
		logger:            logger,
	}
}

func (s *BookingService) CreateBooking(userID uuid.UUID, req *model.CreateBookingRequest) (*model.BookingResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}

	// Create booking model
	booking := &model.Booking{
		UserID:          userID,
		ExpertID:        req.ExpertID,
		ScheduledTime:   req.ScheduledTime,
		DurationMinutes: req.DurationMinutes,
		MeetingType:     model.BookingType(req.MeetingType),
		Status:          model.BookingStatusPending,
		MeetingAddress:  derefString(req.MeetingAddress),
		MeetingURL:      derefString(req.MeetingURL),
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save booking to database
	createdBooking, err := s.bookingRepo.Create(booking)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %v", err)
	}

	// Create initial status history
	statusHistory := &model.StatusHistory{
		BookingID: createdBooking.ID,
		OldStatus: "",
		NewStatus: model.BookingStatusPending,
		ChangedBy: userID,
		Reason:    "Booking created",
		CreatedAt: time.Now(),
	}

	if err := s.statusHistoryRepo.Create(statusHistory); err != nil {
		s.logger.Error("Failed to create status history", err)
	}

	// Cache booking data
	s.cacheBooking(createdBooking)

	// Clear expert and user availability cache
	ctx := context.Background()
	expertCachePattern := fmt.Sprintf("expert_*:%s:*", createdBooking.ExpertID.String())
	userCachePattern := fmt.Sprintf("user_*:%s:*", createdBooking.UserID.String())

	// Delete all expert related cache
	if keys, err := s.redisClient.Keys(ctx, expertCachePattern).Result(); err == nil && len(keys) > 0 {
		s.redisClient.Del(ctx, keys...)
	}

	// Delete all user related cache
	if keys, err := s.redisClient.Keys(ctx, userCachePattern).Result(); err == nil && len(keys) > 0 {
		s.redisClient.Del(ctx, keys...)
	}

	// TODO: Send notification to expert
	s.notifyBookingCreated(createdBooking)

	return s.convertToBookingResponse(createdBooking), nil
}

func (s *BookingService) GetBookingByID(bookingID uuid.UUID) (*model.BookingResponse, error) {
	// Try to get from cache first
	if booking := s.getBookingFromCache(bookingID); booking != nil {
		return s.convertToBookingResponse(booking), nil
	}

	// Get from database
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.cacheBooking(booking)

	return s.convertToBookingResponse(booking), nil
}

func (s *BookingService) UpdateBooking(bookingID uuid.UUID, req *model.UpdateBookingRequest) (*model.BookingResponse, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return nil, err
	}

	// Check if booking can be updated based on status
	if booking.Status == model.BookingStatusCancelled {
		return nil, fmt.Errorf("cannot update cancelled booking")
	}
	if booking.Status == model.BookingStatusCompleted {
		return nil, fmt.Errorf("cannot update completed booking")
	}
	if booking.Status == model.BookingStatusRejected {
		return nil, fmt.Errorf("cannot update rejected booking")
	}

	// Update fields
	if req.ScheduledTime != nil {
		booking.ScheduledTime = *req.ScheduledTime
	}
	if req.DurationMinutes != nil {
		booking.DurationMinutes = *req.DurationMinutes
	}
	if req.Notes != nil {
		booking.Notes = *req.Notes
	}
	if req.MeetingType != nil {
		booking.MeetingType = model.BookingType(*req.MeetingType)
	}
	if req.MeetingAddress != nil {
		booking.MeetingAddress = *req.MeetingAddress
	}
	if req.MeetingURL != nil {
		booking.MeetingURL = *req.MeetingURL
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

	return s.convertToBookingResponse(updatedBooking), nil
}

func (s *BookingService) CancelBooking(bookingID uuid.UUID, userID uuid.UUID, req *model.CancelBookingRequest) error {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return err
	}

	// Check if booking can be cancelled
	if !booking.CanBeCancelled() {
		return fmt.Errorf("booking cannot be cancelled")
	}

	// Update status to cancelled
	err = s.bookingRepo.UpdateStatus(bookingID, model.BookingStatusCancelled)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %v", err)
	}

	// Create status history
	statusHistory := &model.StatusHistory{
		BookingID: bookingID,
		OldStatus: booking.Status,
		NewStatus: model.BookingStatusCancelled,
		ChangedBy: userID,
		Reason:    req.Reason,
		CreatedAt: time.Now(),
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

func (s *BookingService) GetUserBookings(userID uuid.UUID, req *model.GetBookingsRequest) ([]model.BookingResponse, int64, error) {
	filter := &model.BookingFilter{
		UserID:    &userID,
		Status:    req.Status,
		Type:      req.Type,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Page:      req.Page,
		Limit:     req.Limit,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	bookings, total, err := s.bookingRepo.GetByUserID(userID, filter)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	responses := make([]model.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = *s.convertToBookingResponse(&booking)
	}

	return responses, total, nil
}

func (s *BookingService) GetExpertBookings(expertID uuid.UUID, req *model.GetBookingsRequest) ([]model.BookingResponse, int64, error) {
	filter := &model.BookingFilter{
		ExpertID:  &expertID,
		Status:    req.Status,
		Type:      req.Type,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Page:      req.Page,
		Limit:     req.Limit,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	bookings, total, err := s.bookingRepo.GetByExpertID(expertID, filter)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	responses := make([]model.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = *s.convertToBookingResponse(&booking)
	}

	return responses, total, nil
}

func (s *BookingService) GetBookingHistory(userID uuid.UUID, req *model.GetHistoryRequest) ([]model.StatusHistoryResponse, int64, error) {
	histories, total, err := s.statusHistoryRepo.GetHistoryByUserID(userID, req)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	responses := make([]model.StatusHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = model.StatusHistoryResponse{
			StatusHistory: &history,
		}
	}

	return responses, total, nil
}

func (s *BookingService) GetExpertHistory(expertID uuid.UUID, req *model.GetHistoryRequest) ([]model.StatusHistoryResponse, int64, error) {
	histories, total, err := s.statusHistoryRepo.GetHistoryByExpertID(expertID, req)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	responses := make([]model.StatusHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = model.StatusHistoryResponse{
			StatusHistory: &history,
		}
	}

	return responses, total, nil
}

func (s *BookingService) GetBookingStatistics(userID uuid.UUID, userRole, period string, year, month int) (*model.BookingStatsResponse, error) {
	if userRole == "expert" {
		return s.bookingRepo.GetBookingStats(nil, &userID, nil, nil)
	}
	return s.bookingRepo.GetBookingStats(&userID, nil, nil, nil)
}

func (s *BookingService) GetUpcomingUserBookings(userID uuid.UUID, limit, days int) ([]model.BookingResponse, error) {
	endDate := time.Now().AddDate(0, 0, days)
	bookings, err := s.bookingRepo.GetUpcomingByUserID(userID, limit, endDate)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]model.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = *s.convertToBookingResponse(booking)
	}

	return responses, nil
}

func (s *BookingService) GetUpcomingExpertBookings(expertID uuid.UUID, limit, days int) ([]model.BookingResponse, error) {
	endDate := time.Now().AddDate(0, 0, days)
	bookings, err := s.bookingRepo.GetUpcomingByExpertID(expertID, limit, endDate)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]model.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = *s.convertToBookingResponse(booking)
	}

	return responses, nil
}

func (s *BookingService) GetPastUserBookings(userID uuid.UUID, page, limit int) ([]model.BookingResponse, int64, error) {
	bookings, total, err := s.bookingRepo.GetPastByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	responses := make([]model.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = *s.convertToBookingResponse(booking)
	}

	return responses, int64(total), nil
}

func (s *BookingService) GetPastExpertBookings(expertID uuid.UUID, page, limit int) ([]model.BookingResponse, int64, error) {
	bookings, total, err := s.bookingRepo.GetPastByExpertID(expertID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	responses := make([]model.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = *s.convertToBookingResponse(booking)
	}

	return responses, int64(total), nil
}

func (s *BookingService) CheckConflict(req *model.CheckConflictRequest) (*model.ConflictCheckResponse, error) {
	conflicts, err := s.bookingRepo.CheckConflict(req)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]model.BookingResponse, len(conflicts))
	for i, booking := range conflicts {
		responses[i] = *s.convertToBookingResponse(&booking)
	}

	return &model.ConflictCheckResponse{
		HasConflict:      len(conflicts) > 0,
		ConflictBookings: responses,
	}, nil
}

func (s *BookingService) CheckConflictWithExclusion(req *model.CheckConflictRequest, excludeID uuid.UUID) (*model.ConflictCheckResponse, error) {
	req.ExcludeID = &excludeID
	conflicts, err := s.bookingRepo.CheckConflict(req)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]model.BookingResponse, len(conflicts))
	for i, booking := range conflicts {
		responses[i] = *s.convertToBookingResponse(&booking)
	}

	return &model.ConflictCheckResponse{
		HasConflict:      len(conflicts) > 0,
		ConflictBookings: responses,
	}, nil
}

func (s *BookingService) GetExpertBookingsByDate(expertID uuid.UUID, date time.Time) ([]model.BookingResponse, error) {
	bookings, err := s.bookingRepo.GetExpertBookingsByDate(expertID, date)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]model.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = *s.convertToBookingResponse(&booking)
	}

	return responses, nil
}

func (s *BookingService) GetExpertIDByUserID(userID uuid.UUID) (uuid.UUID, error) {
	return s.bookingRepo.GetExpertIDByUserID(userID)
}

// DeleteBooking deletes a booking and clears related cache
func (s *BookingService) DeleteBooking(bookingID uuid.UUID) error {
	// Get booking before deletion to get expert and user IDs
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %v", err)
	}

	// Delete from database
	if err := s.bookingRepo.Delete(bookingID); err != nil {
		return fmt.Errorf("failed to delete booking: %v", err)
	}

	// Clear all related cache
	ctx := context.Background()

	// 1. Clear booking cache
	bookingKey := fmt.Sprintf("booking:%s", bookingID.String())
	s.redisClient.Del(ctx, bookingKey)

	// 2. Clear expert cache
	expertCachePattern := fmt.Sprintf("expert_*:%s:*", booking.ExpertID.String())
	if keys, err := s.redisClient.Keys(ctx, expertCachePattern).Result(); err == nil && len(keys) > 0 {
		s.redisClient.Del(ctx, keys...)
	}

	// 3. Clear user cache
	userCachePattern := fmt.Sprintf("user_*:%s:*", booking.UserID.String())
	if keys, err := s.redisClient.Keys(ctx, userCachePattern).Result(); err == nil && len(keys) > 0 {
		s.redisClient.Del(ctx, keys...)
	}

	// 4. Clear expert busy time cache
	busyTimePattern := fmt.Sprintf("expert_busy:%s:*", booking.ExpertID.String())
	if keys, err := s.redisClient.Keys(ctx, busyTimePattern).Result(); err == nil && len(keys) > 0 {
		s.redisClient.Del(ctx, keys...)
	}

	return nil
}

// Helper function to convert Booking to BookingResponse
func (s *BookingService) convertToBookingResponse(booking *model.Booking) *model.BookingResponse {
	return &model.BookingResponse{
		Booking:        booking,
		CanBeCancelled: booking.CanBeCancelled(),
		CanBeConfirmed: booking.CanBeConfirmed(),
		IsExpired:      booking.IsExpired(),
		Duration:       booking.DurationMinutes,
	}
}

// Helper function to cache booking data
func (s *BookingService) cacheBooking(booking *model.Booking) {
	ctx := context.Background()
	key := fmt.Sprintf("booking:%s", booking.ID.String())
	data, err := json.Marshal(booking)
	if err != nil {
		s.logger.Error("Failed to marshal booking data", err)
		return
	}

	err = s.redisClient.Set(ctx, key, data, 24*time.Hour).Err()
	if err != nil {
		s.logger.Error("Failed to cache booking data", err)
	}
}

// Helper function to get booking from cache
func (s *BookingService) getBookingFromCache(bookingID uuid.UUID) *model.Booking {
	ctx := context.Background()
	key := fmt.Sprintf("booking:%s", bookingID.String())
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil
	}

	var booking model.Booking
	if err := json.Unmarshal(data, &booking); err != nil {
		s.logger.Error("Failed to unmarshal booking data", err)
		return nil
	}

	return &booking
}

// Helper function to notify about booking creation
func (s *BookingService) notifyBookingCreated(booking *model.Booking) {
	// TODO: Implement notification logic
}

// Helper function to notify about booking update
func (s *BookingService) notifyBookingUpdated(booking *model.Booking) {
	// TODO: Implement notification logic
}

// Helper function to notify about booking cancellation
func (s *BookingService) notifyBookingCancelled(booking *model.Booking) {
	// TODO: Implement notification logic
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

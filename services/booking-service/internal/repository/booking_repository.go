// services/booking-service/internal/repository/booking_repository.go - Thịnh
package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"services/booking-service/internal/model"
)

// BookingRepositoryInterface defines the interface for booking repository operations
type BookingRepositoryInterface interface {
	// Basic CRUD operations
	Create(booking *model.Booking) (*model.Booking, error)
	GetByID(id uuid.UUID) (*model.Booking, error)
	Update(booking *model.Booking) (*model.Booking, error)
	Delete(id uuid.UUID) error

	// List and filter operations
	List(filter *model.BookingFilter) ([]model.Booking, int64, error)
	GetByUserID(userID uuid.UUID, filter *model.BookingFilter) ([]model.Booking, int64, error)
	GetByExpertID(expertID uuid.UUID, filter *model.BookingFilter) ([]model.Booking, int64, error)
	GetBookingsByDateRange(startDate, endDate time.Time) ([]model.Booking, error)

	// Status operations
	UpdateStatus(id uuid.UUID, status model.BookingStatus) error
	GetActiveBookingsByExpert(expertID uuid.UUID) ([]model.Booking, error)
	GetActiveBookingsByUser(userID uuid.UUID) ([]model.Booking, error)

	// Conflict checking
	CheckConflict(req *model.CheckConflictRequest) ([]model.Booking, error)
	HasExpertConflict(expertID uuid.UUID, startTime, endTime time.Time) (bool, error)
	HasUserConflict(userID uuid.UUID, startTime, endTime time.Time) (bool, error)
	GetExpertBookingsByDate(expertID uuid.UUID, date time.Time) ([]model.Booking, error)

	// History and statistics
	GetHistoryByUserID(userID uuid.UUID, offset, limit int, status string, startDate, endDate *time.Time) ([]*model.Booking, int, error)
	GetHistoryByExpertID(expertID uuid.UUID, offset, limit int, status string, startDate, endDate *time.Time) ([]*model.Booking, int, error)
	GetUpcomingByUserID(userID uuid.UUID, limit int, endDate time.Time) ([]*model.Booking, error)
	GetUpcomingByExpertID(expertID uuid.UUID, limit int, endDate time.Time) ([]*model.Booking, error)
	GetPastByUserID(userID uuid.UUID, offset, limit int) ([]*model.Booking, int, error)
	GetPastByExpertID(expertID uuid.UUID, offset, limit int) ([]*model.Booking, int, error)
	GetUserStatistics(userID uuid.UUID, period string, year, month int) (map[string]interface{}, error)
	GetExpertStatistics(expertID uuid.UUID, period string, year, month int) (map[string]interface{}, error)
	GetBookingStats(userID *uuid.UUID, expertID *uuid.UUID, startDate, endDate *time.Time) (*model.BookingStatsResponse, error)

	// System operations
	GetUpcomingBookings(minutes int) ([]model.Booking, error)
	GetExpiredBookings() ([]model.Booking, error)
}

// bookingRepository struct implement BookingRepositoryInterface
type bookingRepository struct {
	db *gorm.DB
}

// NewBookingRepository creates a new instance of BookingRepository
func NewBookingRepository(db *gorm.DB) BookingRepositoryInterface {
	return &bookingRepository{
		db: db,
	}
}

// Create creates a new booking
func (r *bookingRepository) Create(booking *model.Booking) (*model.Booking, error) {
	err := r.db.Create(booking).Error
	return booking, err
}

// GetByID gets a booking by ID
func (r *bookingRepository) GetByID(id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// Update updates a booking
func (r *bookingRepository) Update(booking *model.Booking) (*model.Booking, error) {
	err := r.db.Save(booking).Error
	return booking, err
}

// Delete deletes a booking
func (r *bookingRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Booking{}, "id = ?", id).Error
}

// List gets a list of bookings with filters
func (r *bookingRepository) List(filter *model.BookingFilter) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var total int64

	query := r.db.Model(&model.Booking{})

	// Apply filters
	query = r.applyFilters(query, filter)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	query = r.applySortingAndPagination(query, filter)

	err := query.Find(&bookings).Error
	return bookings, total, err
}

// GetByUserID gets bookings by user ID
func (r *bookingRepository) GetByUserID(userID uuid.UUID, filter *model.BookingFilter) ([]model.Booking, int64, error) {
	if filter == nil {
		filter = &model.BookingFilter{}
	}
	filter.UserID = &userID
	return r.List(filter)
}

// GetByExpertID gets bookings by expert ID
func (r *bookingRepository) GetByExpertID(expertID uuid.UUID, filter *model.BookingFilter) ([]model.Booking, int64, error) {
	if filter == nil {
		filter = &model.BookingFilter{}
	}
	filter.ExpertID = &expertID
	return r.List(filter)
}

// CheckConflict checks for booking time conflicts
func (r *bookingRepository) CheckConflict(req *model.CheckConflictRequest) ([]model.Booking, error) {
	var conflictBookings []model.Booking

	query := r.db.Where("(expert_id = ? OR user_id = ?) AND status IN (?, ?)",
		req.ExpertID, req.UserID, model.BookingStatusPending, model.BookingStatusConfirmed)

	// Check time overlap
	query = query.Where("(scheduled_datetime < ? AND scheduled_datetime + (duration_minutes || ' minutes')::interval > ?) OR "+
		"(scheduled_datetime < ? AND scheduled_datetime + (duration_minutes || ' minutes')::interval > ?) OR "+
		"(scheduled_datetime >= ? AND scheduled_datetime < ?)",
		req.EndTime, req.StartTime, req.StartTime, req.EndTime, req.StartTime, req.EndTime)

	// Exclude current booking if any
	if req.ExcludeID != nil {
		query = query.Where("id != ?", *req.ExcludeID)
	}

	err := query.Find(&conflictBookings).Error
	return conflictBookings, err
}

// HasExpertConflict checks for expert time conflicts
func (r *bookingRepository) HasExpertConflict(expertID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&model.Booking{}).
		Where("expert_id = ? AND status IN (?, ?) AND "+
			"((scheduled_datetime < ? AND scheduled_datetime + (duration_minutes || ' minutes')::interval > ?) OR "+
			"(scheduled_datetime < ? AND scheduled_datetime + (duration_minutes || ' minutes')::interval > ?) OR "+
			"(scheduled_datetime >= ? AND scheduled_datetime < ?))",
			expertID, model.BookingStatusPending, model.BookingStatusConfirmed,
			endTime, startTime, startTime, endTime, startTime, endTime).
		Count(&count).Error
	return count > 0, err
}

// HasUserConflict checks for user time conflicts
func (r *bookingRepository) HasUserConflict(userID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&model.Booking{}).
		Where("user_id = ? AND status IN (?, ?) AND "+
			"((scheduled_datetime < ? AND scheduled_datetime + (duration_minutes || ' minutes')::interval > ?) OR "+
			"(scheduled_datetime < ? AND scheduled_datetime + (duration_minutes || ' minutes')::interval > ?) OR "+
			"(scheduled_datetime >= ? AND scheduled_datetime < ?))",
			userID, model.BookingStatusPending, model.BookingStatusConfirmed,
			endTime, startTime, startTime, endTime, startTime, endTime).
		Count(&count).Error
	return count > 0, err
}

// GetExpertBookingsByDate gets expert bookings for a specific date
func (r *bookingRepository) GetExpertBookingsByDate(expertID uuid.UUID, date time.Time) ([]model.Booking, error) {
	var bookings []model.Booking
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Where("expert_id = ? AND scheduled_datetime >= ? AND scheduled_datetime < ?",
		expertID, startOfDay, endOfDay).
		Order("scheduled_datetime ASC").
		Find(&bookings).Error

	return bookings, err
}

// GetUpcomingBookings gets bookings scheduled in the next X minutes
func (r *bookingRepository) GetUpcomingBookings(minutes int) ([]model.Booking, error) {
	var bookings []model.Booking

	now := time.Now()
	targetTime := now.Add(time.Duration(minutes) * time.Minute)

	err := r.db.Where("status IN (?, ?) AND scheduled_datetime BETWEEN ? AND ?",
		model.BookingStatusPending, model.BookingStatusConfirmed, now, targetTime).
		Find(&bookings).Error

	return bookings, err
}

// GetExpiredBookings gets expired pending bookings
func (r *bookingRepository) GetExpiredBookings() ([]model.Booking, error) {
	var bookings []model.Booking

	now := time.Now()
	err := r.db.Where("status = ? AND scheduled_datetime + (duration_minutes || ' minutes')::interval < ?",
		model.BookingStatusPending, now).
		Find(&bookings).Error

	return bookings, err
}

// UpdateStatus updates booking status
func (r *bookingRepository) UpdateStatus(id uuid.UUID, status model.BookingStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// Update timestamp based on status
	switch status {
	case model.BookingStatusConfirmed:
		updates["confirmed_at"] = time.Now()
	case model.BookingStatusCancelled:
		updates["cancelled_at"] = time.Now()
	case model.BookingStatusCompleted:
		updates["completed_at"] = time.Now()
	}

	return r.db.Model(&model.Booking{}).Where("id = ?", id).Updates(updates).Error
}

// GetBookingsByDateRange gets bookings within a date range
func (r *bookingRepository) GetBookingsByDateRange(startDate, endDate time.Time) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Where("scheduled_datetime BETWEEN ? AND ?", startDate, endDate).
		Order("scheduled_datetime ASC").
		Find(&bookings).Error
	return bookings, err
}

// GetActiveBookingsByExpert gets active bookings for an expert
func (r *bookingRepository) GetActiveBookingsByExpert(expertID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Where("expert_id = ? AND status IN (?, ?)",
		expertID, model.BookingStatusPending, model.BookingStatusConfirmed).
		Order("scheduled_datetime ASC").
		Find(&bookings).Error
	return bookings, err
}

// GetActiveBookingsByUser gets active bookings for a user
func (r *bookingRepository) GetActiveBookingsByUser(userID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Where("user_id = ? AND status IN (?, ?)",
		userID, model.BookingStatusPending, model.BookingStatusConfirmed).
		Order("scheduled_datetime ASC").
		Find(&bookings).Error
	return bookings, err
}

// applyFilters applies filters to the query
func (r *bookingRepository) applyFilters(query *gorm.DB, filter *model.BookingFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.ExpertID != nil {
		query = query.Where("expert_id = ?", *filter.ExpertID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Type != nil {
		query = query.Where("meeting_type = ?", *filter.Type)
	}
	if filter.StartDate != nil {
		query = query.Where("scheduled_datetime >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("scheduled_datetime <= ?", *filter.EndDate)
	}

	return query
}

// applySortingAndPagination applies sorting and pagination to the query
func (r *bookingRepository) applySortingAndPagination(query *gorm.DB, filter *model.BookingFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	// Apply sorting
	if filter.SortBy != "" {
		order := "ASC"
		if filter.SortOrder == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", filter.SortBy, order))
	} else {
		// Default sorting
		query = query.Order("scheduled_datetime DESC")
	}

	// Apply pagination
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	return query
}

// GetBookingStats gets booking statistics
func (r *bookingRepository) GetBookingStats(userID *uuid.UUID, expertID *uuid.UUID, startDate, endDate *time.Time) (*model.BookingStatsResponse, error) {
	stats := &model.BookingStatsResponse{
		StatusBreakdown: make(map[string]int64),
		TypeBreakdown:   make(map[string]int64),
	}

	query := r.db.Model(&model.Booking{})

	// Apply filters
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if expertID != nil {
		query = query.Where("expert_id = ?", *expertID)
	}
	if startDate != nil {
		query = query.Where("scheduled_datetime >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("scheduled_datetime <= ?", *endDate)
	}

	// Total bookings
	query.Count(&stats.TotalBookings)

	// Status breakdown
	var statusCounts []struct {
		Status string
		Count  int64
	}

	query.Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts)

	for _, sc := range statusCounts {
		stats.StatusBreakdown[sc.Status] = sc.Count
	}

	// Type breakdown
	var typeCounts []struct {
		Type  string
		Count int64
	}

	query.Select("meeting_type, COUNT(*) as count").
		Group("meeting_type").
		Scan(&typeCounts)

	for _, tc := range typeCounts {
		stats.TypeBreakdown[tc.Type] = tc.Count
	}

	return stats, nil
}

// GetExpertStatistics lấy thống kê booking của expert
func (r *bookingRepository) GetExpertStatistics(expertID uuid.UUID, period string, year, month int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	var startDate, endDate time.Time

	// Xác định khoảng thời gian
	now := time.Now()
	switch period {
	case "day":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)
	case "year":
		startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(1, 0, 0)
	default:
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	// Tổng số booking
	var totalBookings int64
	err := r.db.Model(&model.Booking{}).
		Where("expert_id = ? AND scheduled_datetime BETWEEN ? AND ?", expertID, startDate, endDate).
		Count(&totalBookings).Error
	if err != nil {
		return nil, err
	}
	stats["total_bookings"] = totalBookings

	// Phân loại theo trạng thái
	var statusCounts []struct {
		Status string
		Count  int64
	}
	err = r.db.Model(&model.Booking{}).
		Select("status, COUNT(*) as count").
		Where("expert_id = ? AND scheduled_datetime BETWEEN ? AND ?", expertID, startDate, endDate).
		Group("status").
		Scan(&statusCounts).Error
	if err != nil {
		return nil, err
	}

	statusBreakdown := make(map[string]int64)
	for _, sc := range statusCounts {
		statusBreakdown[sc.Status] = sc.Count
	}
	stats["status_breakdown"] = statusBreakdown

	// Phân loại theo loại booking
	var typeCounts []struct {
		Type  string
		Count int64
	}
	err = r.db.Model(&model.Booking{}).
		Select("meeting_type, COUNT(*) as count").
		Where("expert_id = ? AND scheduled_datetime BETWEEN ? AND ?", expertID, startDate, endDate).
		Group("meeting_type").
		Scan(&typeCounts).Error
	if err != nil {
		return nil, err
	}

	typeBreakdown := make(map[string]int64)
	for _, tc := range typeCounts {
		typeBreakdown[tc.Type] = tc.Count
	}
	stats["type_breakdown"] = typeBreakdown

	// Tỷ lệ hoàn thành
	var completedCount int64
	err = r.db.Model(&model.Booking{}).
		Where("expert_id = ? AND scheduled_datetime BETWEEN ? AND ? AND status = ?",
			expertID, startDate, endDate, model.BookingStatusCompleted).
		Count(&completedCount).Error
	if err != nil {
		return nil, err
	}

	if totalBookings > 0 {
		stats["completion_rate"] = float64(completedCount) / float64(totalBookings)
	} else {
		stats["completion_rate"] = 0.0
	}

	return stats, nil
}

// GetHistoryByExpertID lấy lịch sử booking của expert
func (r *bookingRepository) GetHistoryByExpertID(expertID uuid.UUID, offset, limit int, status string, startDate, endDate *time.Time) ([]*model.Booking, int, error) {
	var bookings []*model.Booking
	var total int64

	query := r.db.Model(&model.Booking{}).Where("expert_id = ?", expertID)

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != nil {
		query = query.Where("scheduled_datetime >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("scheduled_datetime <= ?", *endDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	// Get bookings
	err := query.Order("scheduled_datetime DESC").Find(&bookings).Error
	if err != nil {
		return nil, 0, err
	}

	return bookings, int(total), nil
}

// GetHistoryByUserID lấy lịch sử booking của user
func (r *bookingRepository) GetHistoryByUserID(userID uuid.UUID, offset, limit int, status string, startDate, endDate *time.Time) ([]*model.Booking, int, error) {
	var bookings []*model.Booking
	var total int64

	query := r.db.Model(&model.Booking{}).Where("user_id = ?", userID)

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != nil {
		query = query.Where("scheduled_datetime >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("scheduled_datetime <= ?", *endDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	// Get bookings
	err := query.Order("scheduled_datetime DESC").Find(&bookings).Error
	if err != nil {
		return nil, 0, err
	}

	return bookings, int(total), nil
}

// GetPastByExpertID lấy danh sách booking đã qua của expert
func (r *bookingRepository) GetPastByExpertID(expertID uuid.UUID, offset, limit int) ([]*model.Booking, int, error) {
	var bookings []*model.Booking
	var total int64

	query := r.db.Model(&model.Booking{}).
		Where("expert_id = ? AND scheduled_datetime + (duration_minutes || ' minutes')::interval < ?", expertID, time.Now())

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	// Get bookings
	err := query.Order("scheduled_datetime DESC").Find(&bookings).Error
	if err != nil {
		return nil, 0, err
	}

	return bookings, int(total), nil
}

// GetPastByUserID lấy danh sách booking đã qua của user
func (r *bookingRepository) GetPastByUserID(userID uuid.UUID, offset, limit int) ([]*model.Booking, int, error) {
	var bookings []*model.Booking
	var total int64

	query := r.db.Model(&model.Booking{}).
		Where("user_id = ? AND scheduled_datetime + (duration_minutes || ' minutes')::interval < ?", userID, time.Now())

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	// Get bookings
	err := query.Order("scheduled_datetime DESC").Find(&bookings).Error
	if err != nil {
		return nil, 0, err
	}

	return bookings, int(total), nil
}

// GetUpcomingByExpertID lấy danh sách booking sắp tới của expert
func (r *bookingRepository) GetUpcomingByExpertID(expertID uuid.UUID, limit int, endDate time.Time) ([]*model.Booking, error) {
	var bookings []*model.Booking

	query := r.db.Model(&model.Booking{}).
		Where("expert_id = ? AND scheduled_datetime > ? AND scheduled_datetime <= ? AND status IN (?, ?)",
			expertID, time.Now(), endDate, model.BookingStatusPending, model.BookingStatusConfirmed)

	// Apply limit
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Get bookings
	err := query.Order("scheduled_datetime ASC").Find(&bookings).Error
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

// GetUpcomingByUserID lấy danh sách booking sắp tới của user
func (r *bookingRepository) GetUpcomingByUserID(userID uuid.UUID, limit int, endDate time.Time) ([]*model.Booking, error) {
	var bookings []*model.Booking

	query := r.db.Model(&model.Booking{}).
		Where("user_id = ? AND scheduled_datetime > ? AND scheduled_datetime <= ? AND status IN (?, ?)",
			userID, time.Now(), endDate, model.BookingStatusPending, model.BookingStatusConfirmed)

	// Apply limit
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Get bookings
	err := query.Order("scheduled_datetime ASC").Find(&bookings).Error
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

// GetUserStatistics lấy thống kê booking của user
func (r *bookingRepository) GetUserStatistics(userID uuid.UUID, period string, year, month int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	var startDate, endDate time.Time

	// Xác định khoảng thời gian
	now := time.Now()
	switch period {
	case "day":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)
	case "year":
		startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(1, 0, 0)
	default:
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	// Tổng số booking
	var totalBookings int64
	err := r.db.Model(&model.Booking{}).
		Where("user_id = ? AND scheduled_datetime BETWEEN ? AND ?", userID, startDate, endDate).
		Count(&totalBookings).Error
	if err != nil {
		return nil, err
	}
	stats["total_bookings"] = totalBookings

	// Phân loại theo trạng thái
	var statusCounts []struct {
		Status string
		Count  int64
	}
	err = r.db.Model(&model.Booking{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ? AND scheduled_datetime BETWEEN ? AND ?", userID, startDate, endDate).
		Group("status").
		Scan(&statusCounts).Error
	if err != nil {
		return nil, err
	}

	statusBreakdown := make(map[string]int64)
	for _, sc := range statusCounts {
		statusBreakdown[sc.Status] = sc.Count
	}
	stats["status_breakdown"] = statusBreakdown

	// Phân loại theo loại booking
	var typeCounts []struct {
		Type  string
		Count int64
	}
	err = r.db.Model(&model.Booking{}).
		Select("meeting_type, COUNT(*) as count").
		Where("user_id = ? AND scheduled_datetime BETWEEN ? AND ?", userID, startDate, endDate).
		Group("meeting_type").
		Scan(&typeCounts).Error
	if err != nil {
		return nil, err
	}

	typeBreakdown := make(map[string]int64)
	for _, tc := range typeCounts {
		typeBreakdown[tc.Type] = tc.Count
	}
	stats["type_breakdown"] = typeBreakdown

	// Tỷ lệ hoàn thành
	var completedCount int64
	err = r.db.Model(&model.Booking{}).
		Where("user_id = ? AND scheduled_datetime BETWEEN ? AND ? AND status = ?",
			userID, startDate, endDate, model.BookingStatusCompleted).
		Count(&completedCount).Error
	if err != nil {
		return nil, err
	}

	if totalBookings > 0 {
		stats["completion_rate"] = float64(completedCount) / float64(totalBookings)
	} else {
		stats["completion_rate"] = 0.0
	}

	return stats, nil
}

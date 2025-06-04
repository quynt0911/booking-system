// services/booking-service/internal/repository/booking_repository.go - Thịnh
package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	
	"booking-system/services/booking-service/internal/model"
)

// BookingRepository interface định nghĩa các method cho booking repository
type BookingRepository interface {
	Create(booking *model.Booking) error
	GetByID(id uint) (*model.Booking, error)
	Update(booking *model.Booking) error
	Delete(id uint) error
	List(filter *model.BookingFilter) ([]model.Booking, int64, error)
	GetByUserID(userID uint, filter *model.BookingFilter) ([]model.Booking, int64, error)
	GetByExpertID(expertID uint, filter *model.BookingFilter) ([]model.Booking, int64, error)
	CheckConflict(req *model.CheckConflictRequest) ([]model.Booking, error)
	GetUpcomingBookings(minutes int) ([]model.Booking, error)
	GetExpiredBookings() ([]model.Booking, error)
	GetBookingStats(userID *uint, expertID *uint, startDate, endDate *time.Time) (*model.BookingStatsResponse, error)
	UpdateStatus(id uint, status model.BookingStatus) error
	GetBookingsByDateRange(startDate, endDate time.Time) ([]model.Booking, error)
	GetActiveBookingsByExpert(expertID uint) ([]model.Booking, error)
	GetActiveBookingsByUser(userID uint) ([]model.Booking, error)
}

// bookingRepository struct implement BookingRepository interface
type bookingRepository struct {
	db *gorm.DB
}

// NewBookingRepository tạo instance mới của BookingRepository
func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{
		db: db,
	}
}

// Create tạo booking mới
func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

// GetByID lấy booking theo ID
func (r *bookingRepository) GetByID(id uint) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.First(&booking, id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// Update cập nhật booking
func (r *bookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

// Delete xóa booking (soft delete)
func (r *bookingRepository) Delete(id uint) error {
	return r.db.Delete(&model.Booking{}, id).Error
}

// List lấy danh sách booking với filter
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

// GetByUserID lấy booking theo user ID
func (r *bookingRepository) GetByUserID(userID uint, filter *model.BookingFilter) ([]model.Booking, int64, error) {
	if filter == nil {
		filter = &model.BookingFilter{}
	}
	filter.UserID = &userID
	return r.List(filter)
}

// GetByExpertID lấy booking theo expert ID
func (r *bookingRepository) GetByExpertID(expertID uint, filter *model.BookingFilter) ([]model.Booking, int64, error) {
	if filter == nil {
		filter = &model.BookingFilter{}
	}
	filter.ExpertID = &expertID
	return r.List(filter)
}

// CheckConflict kiểm tra xung đột thời gian booking
func (r *bookingRepository) CheckConflict(req *model.CheckConflictRequest) ([]model.Booking, error) {
	var conflictBookings []model.Booking
	
	query := r.db.Where("(expert_id = ? OR user_id = ?) AND status IN (?, ?)", 
		req.ExpertID, req.UserID, model.StatusPending, model.StatusConfirmed)
	
	// Kiểm tra overlap thời gian
	query = query.Where("(start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?)",
		req.EndTime, req.StartTime, req.StartTime, req.EndTime, req.StartTime, req.EndTime)
	
	// Loại trừ booking hiện tại nếu có
	if req.ExcludeID != nil {
		query = query.Where("id != ?", *req.ExcludeID)
	}
	
	err := query.Find(&conflictBookings).Error
	return conflictBookings, err
}

// GetUpcomingBookings lấy các booking sắp diễn ra trong X phút
func (r *bookingRepository) GetUpcomingBookings(minutes int) ([]model.Booking, error) {
	var bookings []model.Booking
	
	now := time.Now()
	targetTime := now.Add(time.Duration(minutes) * time.Minute)
	
	err := r.db.Where("status IN (?, ?) AND start_time BETWEEN ? AND ?",
		model.StatusPending, model.StatusConfirmed, now, targetTime).
		Find(&bookings).Error
		
	return bookings, err
}

// GetExpiredBookings lấy các booking đã hết hạn nhưng vẫn pending
func (r *bookingRepository) GetExpiredBookings() ([]model.Booking, error) {
	var bookings []model.Booking
	
	now := time.Now()
	err := r.db.Where("status = ? AND end_time < ?", model.StatusPending, now).
		Find(&bookings).Error
		
	return bookings, err
}

// GetBookingStats lấy thống kê booking
func (r *bookingRepository) GetBookingStats(userID *uint, expertID *uint, startDate, endDate *time.Time) (*model.BookingStatsResponse, error) {
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
		query = query.Where("start_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("start_time <= ?", *endDate)
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
		
		switch sc.Status {
		case string(model.StatusPending):
			stats.PendingBookings = sc.Count
		case string(model.StatusConfirmed):
			stats.ConfirmedBookings = sc.Count
		case string(model.StatusCompleted):
			stats.CompletedBookings = sc.Count
		case string(model.StatusCancelled):
			stats.CancelledBookings = sc.Count
		}
	}
	
	// Type breakdown
	var typeCounts []struct {
		Type  string
		Count int64
	}
	
	query.Select("type, COUNT(*) as count").
		Group("type").
		Scan(&typeCounts)
		
	for _, tc := range typeCounts {
		stats.TypeBreakdown[tc.Type] = tc.Count
	}
	
	return stats, nil
}

// UpdateStatus cập nhật trạng thái booking
func (r *bookingRepository) UpdateStatus(id uint, status model.BookingStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	
	// Cập nhật timestamp tương ứng
	switch status {
	case model.StatusConfirmed:
		now := time.Now()
		updates["confirmed_at"] = &now
	case model.StatusCancelled:
		now := time.Now()
		updates["cancelled_at"] = &now
	}
	
	return r.db.Model(&model.Booking{}).Where("id = ?", id).Updates(updates).Error
}

// GetBookingsByDateRange lấy booking theo khoảng thời gian
func (r *bookingRepository) GetBookingsByDateRange(startDate, endDate time.Time) ([]model.Booking, error) {
	var bookings []model.Booking
	
	err := r.db.Where("start_time >= ? AND start_time <= ?", startDate, endDate).
		Order("start_time ASC").
		Find(&bookings).Error
		
	return bookings, err
}

// GetActiveBookingsByExpert lấy booking đang active của expert
func (r *bookingRepository) GetActiveBookingsByExpert(expertID uint) ([]model.Booking, error) {
	var bookings []model.Booking
	
	err := r.db.Where("expert_id = ? AND status IN (?, ?)", 
		expertID, model.StatusPending, model.StatusConfirmed).
		Order("start_time ASC").
		Find(&bookings).Error
		
	return bookings, err
}

// GetActiveBookingsByUser lấy booking đang active của user
func (r *bookingRepository) GetActiveBookingsByUser(userID uint) ([]model.Booking, error) {
	var bookings []model.Booking
	
	err := r.db.Where("user_id = ? AND status IN (?, ?)", 
		userID, model.StatusPending, model.StatusConfirmed).
		Order("start_time ASC").
		Find(&bookings).Error
		
	return bookings, err
}

// applyFilters áp dụng các filter cho query
func (r *bookingRepository) applyFilters(query *gorm.DB, filter *model.BookingFilter) *gorm.DB {
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
		query = query.Where("type = ?", *filter.Type)
	}
	
	if filter.StartDate != nil {
		query = query.Where("start_time >= ?", *filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("start_time <= ?", *filter.EndDate)
	}
	
	return query
}

// applySortingAndPagination áp dụng sorting và pagination
func (r *bookingRepository) applySortingAndPagination(query *gorm.DB, filter *model.BookingFilter) *gorm.DB {
	// Sorting
	sortBy := "created_at"
	sortOrder := "desc"
	
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	
	if filter.SortOrder != "" {
		sortOrder = filter.SortOrder
	}
	
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	
	// Pagination
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}
	
	return query
}
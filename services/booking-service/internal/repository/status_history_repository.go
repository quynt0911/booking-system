// services/booking-service/internal/repository/status_history_repository.go - Thịnh
package repository

import (
	"time"

	"gorm.io/gorm"
	
	"booking-system/services/booking-service/internal/model"
)

// StatusHistoryRepository interface định nghĩa các method cho status history repository
type StatusHistoryRepository interface {
	Create(history *model.StatusHistory) error
	GetByID(id uint) (*model.StatusHistory, error)
	GetByBookingID(bookingID uint) ([]model.StatusHistory, error)
	List(filter *model.StatusHistoryFilter) ([]model.StatusHistory, int64, error)
	GetLatestByBookingID(bookingID uint) (*model.StatusHistory, error)
	GetHistoryByUser(changedBy uint, filter *model.StatusHistoryFilter) ([]model.StatusHistory, int64, error)
	DeleteOldHistory(olderThan time.Time) error
	GetStatusChangeCounts(startDate, endDate *time.Time) (map[string]int64, error)
	GetChangesByDateRange(startDate, endDate time.Time) ([]model.StatusHistory, error)
}

// statusHistoryRepository struct implement StatusHistoryRepository interface
type statusHistoryRepository struct {
	db *gorm.DB
}

// NewStatusHistoryRepository tạo instance mới của StatusHistoryRepository
func NewStatusHistoryRepository(db *gorm.DB) StatusHistoryRepository {
	return &statusHistoryRepository{
		db: db,
	}
}

// Create tạo record lịch sử trạng thái mới
func (r *statusHistoryRepository) Create(history *model.StatusHistory) error {
	return r.db.Create(history).Error
}

// GetByID lấy lịch sử trạng thái theo ID
func (r *statusHistoryRepository) GetByID(id uint) (*model.StatusHistory, error) {
	var history model.StatusHistory
	err := r.db.First(&history, id).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// GetByBookingID lấy tất cả lịch sử trạng thái của một booking
func (r *statusHistoryRepository) GetByBookingID(bookingID uint) ([]model.StatusHistory, error) {
	var histories []model.StatusHistory
	
	err := r.db.Where("booking_id = ?", bookingID).
		Order("created_at DESC").
		Find(&histories).Error
		
	return histories, err
}

// List lấy danh sách lịch sử trạng thái với filter
func (r *statusHistoryRepository) List(filter *model.StatusHistoryFilter) ([]model.StatusHistory, int64, error) {
	var histories []model.StatusHistory
	var total int64

	query := r.db.Model(&model.StatusHistory{})
	
	// Apply filters
	query = r.applyFilters(query, filter)
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	query = r.applyPagination(query, filter)
	
	// Order by created_at DESC
	err := query.Order("created_at DESC").Find(&histories).Error
	return histories, total, err
}

// GetLatestByBookingID lấy lịch sử trạng thái mới nhất của booking
func (r *statusHistoryRepository) GetLatestByBookingID(bookingID uint) (*model.StatusHistory, error) {
	var history model.StatusHistory
	
	err := r.db.Where("booking_id = ?", bookingID).
		Order("created_at DESC").
		First(&history).Error
		
	if err != nil {
		return nil, err
	}
	
	return &history, nil
}

// GetHistoryByUser lấy lịch sử thay đổi trạng thái theo user
func (r *statusHistoryRepository) GetHistoryByUser(changedBy uint, filter *model.StatusHistoryFilter) ([]model.StatusHistory, int64, error) {
	if filter == nil {
		filter = &model.StatusHistoryFilter{}
	}
	filter.ChangedBy = &changedBy
	return r.List(filter)
}

// DeleteOldHistory xóa lịch sử cũ hơn thời gian specified
func (r *statusHistoryRepository) DeleteOldHistory(olderThan time.Time) error {
	return r.db.Where("created_at < ?", olderThan).Delete(&model.StatusHistory{}).Error
}

// GetStatusChangeCounts lấy thống kê số lần thay đổi trạng thái
func (r *statusHistoryRepository) GetStatusChangeCounts(startDate, endDate *time.Time) (map[string]int64, error) {
	query := r.db.Model(&model.StatusHistory{})
	
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}
	
	var results []struct {
		NewStatus string
		Count     int64
	}
	
	err := query.Select("new_status, COUNT(*) as count").
		Group("new_status").
		Scan(&results).Error
		
	if err != nil {
		return nil, err
	}
	
	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.NewStatus] = result.Count
	}
	
	return counts, nil
}

// GetChangesByDateRange lấy lịch sử thay đổi theo khoảng thời gian
func (r *statusHistoryRepository) GetChangesByDateRange(startDate, endDate time.Time) ([]model.StatusHistory, error) {
	var histories []model.StatusHistory
	
	err := r.db.Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC").
		Find(&histories).Error
		
	return histories, err
}

// applyFilters áp dụng các filter cho query
func (r *statusHistoryRepository) applyFilters(query *gorm.DB, filter *model.StatusHistoryFilter) *gorm.DB {
	if filter.BookingID != nil {
		query = query.Where("booking_id = ?", *filter.BookingID)
	}
	
	if filter.ChangedBy != nil {
		query = query.Where("changed_by = ?", *filter.ChangedBy)
	}
	
	if filter.ChangeType != nil {
		query = query.Where("change_type = ?", *filter.ChangeType)
	}
	
	if filter.OldStatus != nil {
		query = query.Where("old_status = ?", *filter.OldStatus)
	}
	
	if filter.NewStatus != nil {
		query = query.Where("new_status = ?", *filter.NewStatus)
	}
	
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	
	return query
}

// applyPagination áp dụng pagination cho query
func (r *statusHistoryRepository) applyPagination(query *gorm.DB, filter *model.StatusHistoryFilter) *gorm.DB {
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}
	
	return query
}

// CreateStatusHistory helper function để tạo lịch sử trạng thái
func (r *statusHistoryRepository) CreateStatusHistory(bookingID uint, oldStatus, newStatus model.BookingStatus, changedBy uint, changeType, reason, notes string) error {
	history := &model.StatusHistory{
		BookingID:  bookingID,
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		ChangedBy:  changedBy,
		ChangeType: changeType,
		Reason:     reason,
		Notes:      notes,
	}
	
	return r.Create(history)
}

// GetStatusTransitionStats lấy thống kê chuyển đổi trạng thái
func (r *statusHistoryRepository) GetStatusTransitionStats(startDate, endDate *time.Time) (map[string]map[string]int64, error) {
	query := r.db.Model(&model.StatusHistory{})
	
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}
	
	var results []struct {
		OldStatus string
		NewStatus string
		Count     int64
	}
	
	err := query.Select("old_status, new_status, COUNT(*) as count").
		Group("old_status, new_status").
		Scan(&results).Error
		
	if err != nil {
		return nil, err
	}
	
	stats := make(map[string]map[string]int64)
	for _, result := range results {
		if stats[result.OldStatus] == nil {
			stats[result.OldStatus] = make(map[string]int64)
		}
		stats[result.OldStatus][result.NewStatus] = result.Count
	}
	
	return stats, nil
}

// GetChangeTypeStats lấy thống kê theo loại thay đổi
func (r *statusHistoryRepository) GetChangeTypeStats(startDate, endDate *time.Time) (map[string]int64, error) {
	query := r.db.Model(&model.StatusHistory{})
	
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}
	
	var results []struct {
		ChangeType string
		Count      int64
	}
	
	err := query.Select("change_type, COUNT(*) as count").
		Group("change_type").
		Scan(&results).Error
		
	if err != nil {
		return nil, err
	}
	
	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.ChangeType] = result.Count
	}
	
	return stats, nil
}

// MostActiveUserResult represents the result for most active users
type MostActiveUserResult struct {
	UserID uint  `json:"user_id"`
	Count  int64 `json:"count"`
}

// GetMostActiveUsers lấy danh sách user thay đổi trạng thái nhiều nhất
func (r *statusHistoryRepository) GetMostActiveUsers(limit int, startDate, endDate *time.Time) ([]MostActiveUserResult, error) {
	query := r.db.Model(&model.StatusHistory{})

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	var results []MostActiveUserResult

	err := query.Select("changed_by as user_id, COUNT(*) as count").
		Group("changed_by").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}
// services/booking-service/internal/repository/status_history_repository.go - Thịnh
package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"booking-system/services/booking-service/internal/model"
)

// StatusHistoryRepositoryInterface defines methods for status history repository
type StatusHistoryRepositoryInterface interface {
	Create(history *model.StatusHistory) error
	GetByID(id uuid.UUID) (*model.StatusHistory, error)
	GetByBookingID(bookingID uuid.UUID) ([]model.StatusHistory, error)
	List(filter *model.StatusHistoryFilter) ([]model.StatusHistory, int64, error)
	GetLatestByBookingID(bookingID uuid.UUID) (*model.StatusHistory, error)
	GetHistoryByUser(changedBy uuid.UUID, filter *model.StatusHistoryFilter) ([]model.StatusHistory, int64, error)
	DeleteOldHistory(olderThan time.Time) error
	GetStatusChangeCounts(startDate, endDate *time.Time) (map[string]int64, error)
	GetChangesByDateRange(startDate, endDate time.Time) ([]model.StatusHistory, error)
	GetByStatus(status string) ([]*model.StatusHistory, error)
	GetByUserID(userID uuid.UUID) ([]model.StatusHistory, error)
	GetByExpertID(expertID uuid.UUID) ([]model.StatusHistory, error)
	CreateStatusHistory(bookingID uuid.UUID, oldStatus, newStatus model.BookingStatus, changedBy uuid.UUID, changeType, reason, notes string) error
	GetStatusTransitionStats(startDate, endDate *time.Time) (map[string]map[string]int64, error)
	GetChangeTypeStats(startDate, endDate *time.Time) (map[string]int64, error)
	GetMostActiveUsers(limit int, startDate, endDate *time.Time) ([]MostActiveUserResult, error)
	GetHistoryByUserID(userID uuid.UUID, req *model.GetHistoryRequest) ([]model.StatusHistory, int64, error)
	GetHistoryByExpertID(expertID uuid.UUID, req *model.GetHistoryRequest) ([]model.StatusHistory, int64, error)
}

// statusHistoryRepository implements StatusHistoryRepositoryInterface
type statusHistoryRepository struct {
	db *gorm.DB
}

// NewStatusHistoryRepository creates a new instance of StatusHistoryRepositoryInterface
func NewStatusHistoryRepository(db *gorm.DB) StatusHistoryRepositoryInterface {
	return &statusHistoryRepository{
		db: db,
	}
}

// Create creates a new status history record
func (r *statusHistoryRepository) Create(history *model.StatusHistory) error {
	return r.db.Create(history).Error
}

// GetByID gets status history by ID
func (r *statusHistoryRepository) GetByID(id uuid.UUID) (*model.StatusHistory, error) {
	var history model.StatusHistory
	err := r.db.First(&history, id).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// GetByBookingID gets all status history for a booking
func (r *statusHistoryRepository) GetByBookingID(bookingID uuid.UUID) ([]model.StatusHistory, error) {
	var histories []model.StatusHistory

	err := r.db.Where("booking_id = ?", bookingID).
		Order("changed_at DESC").
		Find(&histories).Error

	return histories, err
}

// List gets list of status history with filter
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

	// Order by changed_at DESC
	err := query.Order("changed_at DESC").Find(&histories).Error
	return histories, total, err
}

// GetLatestByBookingID gets latest status history for a booking
func (r *statusHistoryRepository) GetLatestByBookingID(bookingID uuid.UUID) (*model.StatusHistory, error) {
	var history model.StatusHistory

	err := r.db.Where("booking_id = ?", bookingID).
		Order("changed_at DESC").
		First(&history).Error

	if err != nil {
		return nil, err
	}

	return &history, nil
}

// GetHistoryByUser gets status change history by user
func (r *statusHistoryRepository) GetHistoryByUser(changedBy uuid.UUID, filter *model.StatusHistoryFilter) ([]model.StatusHistory, int64, error) {
	if filter == nil {
		filter = &model.StatusHistoryFilter{}
	}
	filter.ChangedBy = &changedBy
	return r.List(filter)
}

// DeleteOldHistory deletes history older than specified time
func (r *statusHistoryRepository) DeleteOldHistory(olderThan time.Time) error {
	return r.db.Where("created_at < ?", olderThan).Delete(&model.StatusHistory{}).Error
}

// GetStatusChangeCounts gets statistics of status changes
func (r *statusHistoryRepository) GetStatusChangeCounts(startDate, endDate *time.Time) (map[string]int64, error) {
	query := r.db.Model(&model.StatusHistory{})

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	var results []struct {
		Status string
		Count  int64
	}

	err := query.Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Status] = result.Count
	}

	return counts, nil
}

// GetChangesByDateRange gets history changes by date range
func (r *statusHistoryRepository) GetChangesByDateRange(startDate, endDate time.Time) ([]model.StatusHistory, error) {
	var histories []model.StatusHistory

	err := r.db.Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC").
		Find(&histories).Error

	return histories, err
}

// applyFilters applies filters to query
func (r *statusHistoryRepository) applyFilters(query *gorm.DB, filter *model.StatusHistoryFilter) *gorm.DB {
	if filter.BookingID != nil {
		query = query.Where("booking_id = ?", *filter.BookingID)
	}

	if filter.ChangedBy != nil {
		query = query.Where("changed_by = ?", *filter.ChangedBy)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.StartDate != nil {
		query = query.Where("changed_at >= ?", *filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("changed_at <= ?", *filter.EndDate)
	}

	return query
}

// applyPagination applies pagination to query
func (r *statusHistoryRepository) applyPagination(query *gorm.DB, filter *model.StatusHistoryFilter) *gorm.DB {
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	return query
}

// CreateStatusHistory helper function to create status history
func (r *statusHistoryRepository) CreateStatusHistory(bookingID uuid.UUID, oldStatus, newStatus model.BookingStatus, changedBy uuid.UUID, changeType, reason, notes string) error {
	history := &model.StatusHistory{
		BookingID: bookingID,
		Status:    newStatus,
		ChangedBy: changedBy,
		ChangedAt: time.Now(),
		Note:      notes,
	}

	return r.Create(history)
}

// GetStatusTransitionStats gets status transition statistics
func (r *statusHistoryRepository) GetStatusTransitionStats(startDate, endDate *time.Time) (map[string]map[string]int64, error) {
	query := r.db.Model(&model.StatusHistory{})

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	var results []struct {
		Status string
		Count  int64
	}

	err := query.Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]map[string]int64)
	for _, result := range results {
		if stats["previous"] == nil {
			stats["previous"] = make(map[string]int64)
		}
		stats["previous"][result.Status] = result.Count
	}

	return stats, nil
}

// GetChangeTypeStats gets statistics by change type
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

// MostActiveUserResult represents a user's activity count
type MostActiveUserResult struct {
	UserID uuid.UUID `json:"user_id"`
	Count  int64     `json:"count"`
}

// GetMostActiveUsers gets most active users
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

// GetByStatus gets status history by status
func (r *statusHistoryRepository) GetByStatus(status string) ([]*model.StatusHistory, error) {
	var histories []*model.StatusHistory
	err := r.db.Where("status = ?", status).Find(&histories).Error
	return histories, err
}

// GetByUserID gets status history by user ID
func (r *statusHistoryRepository) GetByUserID(userID uuid.UUID) ([]model.StatusHistory, error) {
	var histories []model.StatusHistory
	err := r.db.Where("changed_by = ?", userID).Find(&histories).Error
	return histories, err
}

// GetByExpertID gets status history by expert ID
func (r *statusHistoryRepository) GetByExpertID(expertID uuid.UUID) ([]model.StatusHistory, error) {
	var histories []model.StatusHistory
	err := r.db.Where("changed_by = ?", expertID).Find(&histories).Error
	return histories, err
}

// GetHistoryByUserID gets status history by user ID with pagination
func (r *statusHistoryRepository) GetHistoryByUserID(userID uuid.UUID, req *model.GetHistoryRequest) ([]model.StatusHistory, int64, error) {
	filter := &model.StatusHistoryFilter{
		ChangedBy: &userID,
		Page:      req.Page,
		Limit:     req.Limit,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	return r.List(filter)
}

// GetHistoryByExpertID gets status history by expert ID with pagination
func (r *statusHistoryRepository) GetHistoryByExpertID(expertID uuid.UUID, req *model.GetHistoryRequest) ([]model.StatusHistory, int64, error) {
	filter := &model.StatusHistoryFilter{
		ChangedBy: &expertID,
		Page:      req.Page,
		Limit:     req.Limit,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	return r.List(filter)
}

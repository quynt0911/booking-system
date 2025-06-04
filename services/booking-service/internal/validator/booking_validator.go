// services/booking-service/internal/validator/booking_validator.go - Thịnh
package validator

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	
	"booking-system/services/booking-service/internal/model"
)

// BookingValidator struct chứa validator instance
type BookingValidator struct {
	validator *validator.Validate
}

// NewBookingValidator tạo instance mới của BookingValidator
func NewBookingValidator() *BookingValidator {
	v := validator.New()
	
	// Đăng ký custom validation functions
	v.RegisterValidation("future_time", validateFutureTime)
	v.RegisterValidation("valid_time_range", validateTimeRange)
	v.RegisterValidation("booking_status", validateBookingStatus)
	v.RegisterValidation("booking_type", validateBookingType)
	v.RegisterValidation("min_duration", validateMinDuration)
	v.RegisterValidation("max_duration", validateMaxDuration)
	v.RegisterValidation("business_hours", validateBusinessHours)
	v.RegisterValidation("not_weekend", validateNotWeekend)
	
	return &BookingValidator{
		validator: v,
	}
}

// ValidateCreateBooking validate request tạo booking
func (bv *BookingValidator) ValidateCreateBooking(req *model.CreateBookingRequest) error {
	if err := bv.validator.Struct(req); err != nil {
		return bv.formatValidationError(err)
	}
	
	// Custom validations
	if err := bv.validateBookingTimeRange(req.StartTime, req.EndTime); err != nil {
		return err
	}
	
	if err := bv.validateBookingType(req.Type, req.Location, req.MeetingLink); err != nil {
		return err
	}
	
	return nil
}

// ValidateUpdateBooking validate request cập nhật booking
func (bv *BookingValidator) ValidateUpdateBooking(req *model.UpdateBookingRequest) error {
	if err := bv.validator.Struct(req); err != nil {
		return bv.formatValidationError(err)
	}
	
	return nil
}

// ValidateStatusUpdate validate request cập nhật trạng thái
func (bv *BookingValidator) ValidateStatusUpdate(req *model.UpdateBookingStatusRequest) error {
	if err := bv.validator.Struct(req); err != nil {
		return bv.formatValidationError(err)
	}
	
	return nil
}

// ValidateConflictCheck validate request kiểm tra xung đột
func (bv *BookingValidator) ValidateConflictCheck(req *model.CheckConflictRequest) error {
	if err := bv.validator.Struct(req); err != nil {
		return bv.formatValidationError(err)
	}
	
	if err := bv.validateBookingTimeRange(req.StartTime, req.EndTime); err != nil {
		return err
	}
	
	return nil
}

// ValidateGetBookings validate request lấy danh sách booking
func (bv *BookingValidator) ValidateGetBookings(req *model.GetBookingsRequest) error {
	if err := bv.validator.Struct(req); err != nil {
		return bv.formatValidationError(err)
	}
	
	// Validate date range
	if req.StartDate != nil && req.EndDate != nil {
		if req.StartDate.After(*req.EndDate) {
			return fmt.Errorf("start_date cannot be after end_date")
		}
	}
	
	// Set default pagination values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}
	
	return nil
}

// validateBookingTimeRange validate thời gian booking
func (bv *BookingValidator) validateBookingTimeRange(startTime, endTime time.Time) error {
	now := time.Now()
	
	// Kiểm tra thời gian trong tương lai
	if startTime.Before(now) {
		return fmt.Errorf("start_time must be in the future")
	}
	
	// Kiểm tra endTime sau startTime
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return fmt.Errorf("end_time must be after start_time")
	}
	
	// Kiểm tra thời gian tối thiểu (15 phút)
	duration := endTime.Sub(startTime)
	if duration < 15*time.Minute {
		return fmt.Errorf("booking duration must be at least 15 minutes")
	}
	
	// Kiểm tra thời gian tối đa (4 giờ)
	if duration > 4*time.Hour {
		return fmt.Errorf("booking duration cannot exceed 4 hours")
	}
	
	// Kiểm tra booking không quá xa trong tương lai (90 ngày)
	maxFutureDate := now.AddDate(0, 0, 90)
	if startTime.After(maxFutureDate) {
		return fmt.Errorf("booking cannot be scheduled more than 90 days in advance")
	}
	
	return nil
}

// validateBookingType validate loại booking và thông tin liên quan
func (bv *BookingValidator) validateBookingType(bookingType model.BookingType, location, meetingLink string) error {
	switch bookingType {
	case model.TypeOnline:
		if location != "" {
			return fmt.Errorf("location should not be provided for online booking")
		}
		// MeetingLink có thể để trống, sẽ được tạo sau
	case model.TypeOffline:
		if meetingLink != "" {
			return fmt.Errorf("meeting_link should not be provided for offline booking")
		}
		if location == "" {
			return fmt.Errorf("location is required for offline booking")
		}
	default:
		return fmt.Errorf("invalid booking type")
	}
	
	return nil
}

// formatValidationError format lỗi validation thành message dễ đọc
func (bv *BookingValidator) formatValidationError(err error) error {
	var errorMessages []string
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errorMessages = append(errorMessages, bv.getFieldErrorMessage(fieldError))
		}
	}
	
	return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, ", "))
}

// getFieldErrorMessage tạo message lỗi cho từng field
func (bv *BookingValidator) getFieldErrorMessage(fieldError validator.FieldError) string {
	field := strings.ToLower(fieldError.Field())
	
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, fieldError.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s", field, fieldError.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fieldError.Param())
	case "future_time":
		return fmt.Sprintf("%s must be in the future", field)
	case "valid_time_range":
		return fmt.Sprintf("invalid time range for %s", field)
	case "booking_status":
		return fmt.Sprintf("invalid booking status: %s", field)
	case "booking_type":
		return fmt.Sprintf("invalid booking type: %s", field)
	case "min_duration":
		return fmt.Sprintf("booking duration must be at least %s minutes", fieldError.Param())
	case "max_duration":
		return fmt.Sprintf("booking duration cannot exceed %s hours", fieldError.Param())
	case "business_hours":
		return fmt.Sprintf("%s must be within business hours (9 AM - 6 PM)", field)
	case "not_weekend":
		return fmt.Sprintf("%s cannot be on weekend", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// Custom validation functions

// validateFutureTime kiểm tra thời gian trong tương lai
func validateFutureTime(fl validator.FieldLevel) bool {
	timeValue, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	
	return timeValue.After(time.Now())
}

// validateTimeRange kiểm tra khoảng thời gian hợp lệ
func validateTimeRange(fl validator.FieldLevel) bool {
	// This would need to be implemented based on struct context
	// For now, return true as we handle this in custom validation
	return true
}

// validateBookingStatus kiểm tra trạng thái booking hợp lệ
func validateBookingStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{
		string(model.StatusPending),
		string(model.StatusConfirmed),
		string(model.StatusRejected),
		string(model.StatusCancelled),
		string(model.StatusCompleted),
		string(model.StatusMissed),
	}
	
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	
	return false
}

// validateBookingType kiểm tra loại booking hợp lệ
func validateBookingType(fl validator.FieldLevel) bool {
	bookingType := fl.Field().String()
	return bookingType == string(model.TypeOnline) || bookingType == string(model.TypeOffline)
}

// validateMinDuration kiểm tra thời gian tối thiểu
func validateMinDuration(fl validator.FieldLevel) bool {
	// This would need access to both start and end time
	// Implemented in custom validation method
	return true
}

// validateMaxDuration kiểm tra thời gian tối đa
func validateMaxDuration(fl validator.FieldLevel) bool {
	// This would need access to both start and end time
	// Implemented in custom validation method
	return true
}

// validateBusinessHours kiểm tra giờ làm việc
func validateBusinessHours(fl validator.FieldLevel) bool {
	timeValue, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	
	hour := timeValue.Hour()
	return hour >= 9 && hour < 18 // 9 AM - 6 PM
}

// validateNotWeekend kiểm tra không phải cuối tuần
func validateNotWeekend(fl validator.FieldLevel) bool {
	timeValue, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	
	weekday := timeValue.Weekday()
	return weekday != time.Saturday && weekday != time.Sunday
}

// ValidateStatusTransition kiểm tra chuyển đổi trạng thái hợp lệ
func (bv *BookingValidator) ValidateStatusTransition(oldStatus, newStatus model.BookingStatus, changeType string) error {
	if !model.IsValidTransition(oldStatus, newStatus, changeType) {
		return fmt.Errorf("invalid status transition from %s to %s by %s", oldStatus, newStatus, changeType)
	}
	
	return nil
}

// ValidateBookingTime kiểm tra thời gian booking với các rule phức tạp
func (bv *BookingValidator) ValidateBookingTime(startTime, endTime time.Time, expertID uint) error {
	// Kiểm tra cơ bản
	if err := bv.validateBookingTimeRange(startTime, endTime); err != nil {
		return err
	}
	
	// Kiểm tra thời gian booking phải là bội số của 15 phút
	startMinute := startTime.Minute()
	endMinute := endTime.Minute()
	
	if startMinute%15 != 0 || endMinute%15 != 0 {
		return fmt.Errorf("booking time must be in 15-minute intervals")
	}
	
	// Kiểm tra không được đặt lịch trong giờ ăn trưa (12:00 - 13:00)
	startHour := startTime.Hour()
	endHour := endTime.Hour()
	
	if (startHour == 12 && startTime.Minute() == 0) || 
	   (endHour == 13 && endTime.Minute() == 0) ||
	   (startHour < 12 && endHour > 13) {
		return fmt.Errorf("booking cannot be scheduled during lunch break (12:00 - 13:00)")
	}
	
	return nil
}
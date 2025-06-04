package utils

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// InitValidator initializes custom validation rules
func InitValidator() {
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("booking_time", validateBookingTime)
	validate.RegisterValidation("consultation_type", validateConsultationType)
}

// ValidateStruct validates a struct using the validator
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// validatePhone validates Vietnamese phone number format
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Vietnamese phone number patterns
	patterns := []string{
		`^(\+84|0)(3[2-9]|5[6|8|9]|7[0|6-9]|8[1-6|8|9]|9[0-4|6-9])[0-9]{7}$`,
		`^(\+84|0)(1[2689]|9[0-9])[0-9]{8}$`,
	}
	
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, phone); matched {
			return true
		}
	}
	return false
}

// validatePassword validates password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	if len(password) < 8 {
		return false
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateBookingTime validates booking time is in the future
func validateBookingTime(fl validator.FieldLevel) bool {
	bookingTime := fl.Field().Interface().(time.Time)
	return bookingTime.After(time.Now())
}

// validateConsultationType validates consultation type
func validateConsultationType(fl validator.FieldLevel) bool {
	consultationType := strings.ToLower(fl.Field().String())
	validTypes := []string{"online", "offline"}
	
	for _, validType := range validTypes {
		if consultationType == validType {
			return true
		}
	}
	return false
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

// GetValidationErrors extracts validation errors
func GetValidationErrors(err error) []ValidationError {
	var errors []ValidationError
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   ve.Field(),
				Tag:     ve.Tag(),
				Message: getErrorMessage(ve),
			})
		}
	}
	
	return errors
}

// getErrorMessage returns user-friendly error message
func getErrorMessage(ve validator.FieldError) string {
	switch ve.Tag() {
	case "required":
		return "Trường này là bắt buộc"
	case "email":
		return "Email không hợp lệ"
	case "phone":
		return "Số điện thoại không hợp lệ"
	case "password":
		return "Mật khẩu phải có ít nhất 8 ký tự, bao gồm chữ hoa, chữ thường, số và ký tự đặc biệt"
	case "min":
		return "Giá trị quá nhỏ"
	case "max":
		return "Giá trị quá lớn"
	case "booking_time":
		return "Thời gian đặt lịch phải trong tương lai"
	case "consultation_type":
		return "Loại tư vấn chỉ có thể là 'online' hoặc 'offline'"
	default:
		return "Giá trị không hợp lệ"
	}
}

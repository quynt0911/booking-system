package utils

// ErrorResponse trả về response lỗi với message và status code
func ErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"message": message,
	}
}

// SuccessResponse trả về response thành công với message và data
func SuccessResponse(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	}
}

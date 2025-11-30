package response

import (
	"github.com/gin-gonic/gin"
)

// Response Standard API (Real Industry Standard)
type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Errors     interface{} `json:"errors"`
}

// Success response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Errors:     nil,
	})
}

// Error response
func Error(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
		Errors:     errors,
	})
}

// ValidationError Helper for validation errors
func ValidationError(c *gin.Context, errors interface{}) {
	Error(c, 400, "Validation failed", errors)
}

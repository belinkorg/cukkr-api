package helper

import (
	"cukurly-app/pkg/errors"
	"cukurly-app/pkg/logger"
	"cukurly-app/pkg/response"
	"cukurly-app/pkg/validator"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPHandlerHelper provide helper methods for HTTP handling
type HTTPHandlerHelper struct {
	validator *validator.Validator
	logger    *logger.Logger
}

// NewHTTPHandlerHelper make instance HTTPHandlerHelper
func NewHTTPHandlerHelper(validator *validator.Validator, logger *logger.Logger) *HTTPHandlerHelper {
	return &HTTPHandlerHelper{
		validator: validator,
		logger:    logger,
	}
}

// HandleError handle error with consistent error response
func (h *HTTPHandlerHelper) HandleError(c *gin.Context, err error, defaultMessage string) {
	if httpErr, ok := err.(*errors.HTTPError); ok {
		h.logger.WithError(err).Warn(httpErr.Message)
		response.Error(c, httpErr.StatusCode, httpErr.Message, err.Error())
	} else {
		h.logger.WithError(err).Error(defaultMessage)
		response.Error(c, http.StatusInternalServerError, defaultMessage, "")
	}
}

// BindAndValidate handle binding JSON and validate request
func (h *HTTPHandlerHelper) BindAndValidate(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.WithError(err).Warn("Invalid request body")
		response.Error(c, http.StatusBadRequest, "Invalid request body", "")
		return err
	}

	if err := h.validator.Validate(req); err != nil {
		h.logger.WithError(err).Warn("Validation failed")
		response.ValidationError(c, err.Error())
		return err
	}

	return nil
}

// ValidateID validate format ID (UUID)
func (h *HTTPHandlerHelper) ValidateID(c *gin.Context, id string) error {
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Invalid ID", "")
		return fmt.Errorf("empty id")
	}

	// Validate UUID format
	if !isValidUUID(id) {
		response.Error(c, http.StatusBadRequest, "Invalid ID format", "")
		return fmt.Errorf("invalid uuid format")
	}

	return nil
}

// ValidateIDWithField validate format ID with custom field name for error message
func (h *HTTPHandlerHelper) ValidateIDWithField(c *gin.Context, id string, fieldName string) error {
	if id == "" {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid %s", fieldName), "")
		return fmt.Errorf("empty %s", fieldName)
	}

	// Validate UUID format
	if !isValidUUID(id) {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid %s format", fieldName), "")
		return fmt.Errorf("invalid %s format", fieldName)
	}

	return nil
}

// isValidUUID checks if string is valid UUID
func isValidUUID(id string) bool {
	_, err := parseUUID(id)
	return err == nil
}

// parseUUID parses UUID string
func parseUUID(id string) (string, error) {
	u := id
	if len(u) != 36 {
		return "", fmt.Errorf("invalid uuid length")
	}
	return u, nil
}

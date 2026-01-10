package cricket

import (
	"errors"
	"net/http"
	"khelogames/logger"
	"github.com/gin-gonic/gin"
)

// Custom error types
var (
	ErrInvalidUUID        = errors.New("invalid UUID format")
	ErrPlayerNotFound     = errors.New("player not found")
	ErrMatchNotFound      = errors.New("match not found")
	ErrTeamNotFound       = errors.New("team not found")
	ErrInningCompleted    = errors.New("inning already completed")
	ErrInvalidBallNumber  = errors.New("invalid ball number")
	ErrInvalidRuns        = errors.New("invalid runs scored")
	ErrBroadcastFailed    = errors.New("failed to broadcast event")
	ErrDatabaseOperation  = errors.New("database operation failed")
)

// Error codes
const (
	CodeInvalidUUID        = "INVALID_UUID"
	CodePlayerNotFound     = "PLAYER_NOT_FOUND"
	CodeMatchNotFound      = "MATCH_NOT_FOUND"
	CodeTeamNotFound       = "TEAM_NOT_FOUND"
	CodeInningCompleted    = "INNING_COMPLETED"
	CodeInvalidBallNumber  = "INVALID_BALL_NUMBER"
	CodeInvalidRuns        = "INVALID_RUNS"
	CodeBroadcastFailed    = "BROADCAST_FAILED"
	CodeDatabaseOperation  = "DATABASE_OPERATION_FAILED"
	CodeValidationFailed   = "VALIDATION_FAILED"
	CodeInternalError      = "INTERNAL_ERROR"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// ErrorHandler handles errors consistently across the cricket package
type ErrorHandler struct {
	logger *logger.Logger
}

func NewErrorHandler(logger *logger.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

// HandleError logs the error and returns appropriate HTTP response
func (eh *ErrorHandler) HandleError(ctx *gin.Context, err error, code string, statusCode int) {
	eh.logger.Error("Cricket operation failed", map[string]interface{}{
		"error": err.Error(),
		"code":  code,
		"path":  ctx.Request.URL.Path,
	})

	response := ErrorResponse{
		Error: err.Error(),
		Code:  code,
	}

	ctx.JSON(statusCode, response)
}

// HandleValidationError handles validation errors
func (eh *ErrorHandler) HandleValidationError(ctx *gin.Context, err error) {
	eh.HandleError(ctx, err, CodeValidationFailed, http.StatusBadRequest)
}

// HandleUUIDError handles UUID parsing errors
func (eh *ErrorHandler) HandleUUIDError(ctx *gin.Context, field string) {
	err := ErrInvalidUUID
	eh.logger.Error("Invalid UUID format", map[string]interface{}{
		"field": field,
		"path":  ctx.Request.URL.Path,
	})
	eh.HandleError(ctx, err, CodeInvalidUUID, http.StatusBadRequest)
}

// HandleNotFoundError handles resource not found errors
func (eh *ErrorHandler) HandleNotFoundError(ctx *gin.Context, resource string) {
	var err error
	var code string
	
	switch resource {
	case "player":
		err = ErrPlayerNotFound
		code = CodePlayerNotFound
	case "match":
		err = ErrMatchNotFound
		code = CodeMatchNotFound
	case "team":
		err = ErrTeamNotFound
		code = CodeTeamNotFound
	default:
		err = errors.New("resource not found")
		code = CodeInternalError
	}
	
	eh.HandleError(ctx, err, code, http.StatusNotFound)
}

// HandleInternalError handles internal server errors
func (eh *ErrorHandler) HandleInternalError(ctx *gin.Context, err error) {
	eh.HandleError(ctx, err, CodeInternalError, http.StatusInternalServerError)
}

// HandleBroadcastError handles broadcast failures
func (eh *ErrorHandler) HandleBroadcastError(ctx *gin.Context, err error) {
	eh.logger.Error("Broadcast failed", map[string]interface{}{
		"error": err.Error(),
		"path":  ctx.Request.URL.Path,
	})
	// Don't return error to client for broadcast failures, just log
}

// ValidateAndParseUUID validates and parses a UUID string
func ValidateAndParseUUID(uuidStr, fieldName string) (string, error) {
	if uuidStr == "" {
		return "", errors.New(fieldName + " is required")
	}
	
	// Basic UUID format validation (36 characters with hyphens)
	if len(uuidStr) != 36 {
		return "", ErrInvalidUUID
	}
	
	return uuidStr, nil
}

// SafeStringConversion safely converts interface{} to string
func SafeStringConversion(value interface{}, fieldName string) (string, error) {
	if value == nil {
		return "", errors.New(fieldName + " is required")
	}
	
	str, ok := value.(string)
	if !ok {
		return "", errors.New(fieldName + " must be a string")
	}
	
	if str == "" {
		return "", errors.New(fieldName + " cannot be empty")
	}
	
	return str, nil
}

// SafeIntConversion safely converts interface{} to int
func SafeIntConversion(value interface{}, fieldName string) (int, error) {
	if value == nil {
		return 0, errors.New(fieldName + " is required")
	}
	
	// Handle both int and float64 (from JSON unmarshaling)
	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	default:
		return 0, errors.New(fieldName + " must be a number")
	}
}




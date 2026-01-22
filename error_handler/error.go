package errorhandler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidationErrorResponse sends a standardized validation error response
// Used for: Invalid input, validation failures, malformed requests
func ValidationErrorResponse(
	ctx *gin.Context,
	fields map[string]string,
) {
	fmt.Println("Fields: ", fields)
	ctx.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid input",
			"fields":  fields,
		},
		"request_id": ctx.GetString("request_id"),
	})
}

// InternalErrorResponse sends a standardized internal server error response
// Used for: Database errors, internal failures, unexpected errors
func InternalErrorResponse(
	ctx *gin.Context,
	message string,
) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": message,
		},
		"request_id": ctx.GetString("request_id"),
	})
}

// NotFoundErrorResponse sends a standardized not found error response
// Used for: Resource not found, invalid IDs
func NotFoundErrorResponse(
	ctx *gin.Context,
	message string,
) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_FOUND",
			"message": message,
		},
		"request_id": ctx.GetString("request_id"),
	})
}

// ForbiddenErrorResponse sends a standardized forbidden error response
// Used for: Unauthorized access, permission denied
func ForbiddenErrorResponse(
	ctx *gin.Context,
	message string,
) {
	ctx.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "FORBIDDEN",
			"message": message,
		},
		"request_id": ctx.GetString("request_id"),
	})
}

// UnauthorizedErrorResponse sends a standardized unauthorized error response
// Used for: Authentication required, invalid token
func UnauthorizedErrorResponse(
	ctx *gin.Context,
	message string,
) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": message,
		},
		"request_id": ctx.GetString("request_id"),
	})
}

// ConflictErrorResponse sends a standardized conflict error response
// Used for: Duplicate entries, constraint violations
func ConflictErrorResponse(
	ctx *gin.Context,
	message string,
) {
	ctx.JSON(http.StatusConflict, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "CONFLICT",
			"message": message,
		},
		"request_id": ctx.GetString("request_id"),
	})
}

// SuccessResponse sends a standardized success response
// Used for: Successful operations with data
func SuccessResponse(
	ctx *gin.Context,
	data interface{},
) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// SuccessResponseWithStatus sends a standardized success response with custom status code
// Used for: Successful operations with specific HTTP status codes (e.g., 201 Created)
func SuccessResponseWithStatus(
	ctx *gin.Context,
	statusCode int,
	data interface{},
) {
	ctx.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

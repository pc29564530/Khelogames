package handlers

import (
	"khelogames/core/token"
	"khelogames/database/models"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"khelogames/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) GetPlayerWithProfileFunc(ctx *gin.Context) {
	var req struct {
		PublicID string `uri:"public_id" binding:"required"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	playerProfile, err := s.store.GetPlayerWithProfile(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get player with profile", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player with profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    playerProfile,
	})
	return
}

func (s *HandlersServer) GetProfileByPublicIDFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get profile")
	var req struct {
		PublicID string `uri:"profile_public_id" binding:"required"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"profile_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	profile, err := s.store.GetProfileByPublicID(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Successfully retrieved profile by profile public_id: ", profile)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

func (s *HandlersServer) GetProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get profile")
	var req struct {
		PublicID string `uri:"public_id" binding:"required"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	profile, err := s.store.GetProfile(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Successfully retrieved profile: ", profile)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

func (s *HandlersServer) GetProfileByFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get profile")
	var req struct {
		PublicID string `uri:"public_id" binding:"required"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	profile, err := s.store.GetProfile(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("Successfully retrieved profile: ", profile)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

type editProfileRequest struct {
	FullName  string `json:"full_name" binding:"required"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url,omitempty"`
	City      string `json:"city" binding:"required"`
	State     string `json:"state" binding:"required"`
	Country   string `json:"country" binding:"required"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func (s *HandlersServer) UpdateProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update profile")
	var req editProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	s.logger.Debug("Request JSON bind successful: ", req)
	var err error
	var emtpyString string
	var latitude float64
	var longitude float64
	if req.Latitude != emtpyString && req.Longitude != emtpyString {
		latitude, err = strconv.ParseFloat(req.Latitude, 64)
		if err != nil {
			s.logger.Error("Failed to parse latitude: ", err)
			fieldErrors := map[string]string{"latitude": "Invalid format"}
			errorhandler.ValidationErrorResponse(ctx, fieldErrors)
			return
		}
		longitude, err = strconv.ParseFloat(req.Longitude, 64)
		if err != nil {
			s.logger.Error("Failed to parse longitude: ", err)
			fieldErrors := map[string]string{"longitude": "Invalid format"}
			errorhandler.ValidationErrorResponse(ctx, fieldErrors)
			return
		}
	}
	var h3Index string
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	update, err := s.txStore.UpdateProfileTx(ctx, authPayload.PublicID, req.Bio, req.AvatarUrl, req.FullName, req.City, req.State, req.Country, latitude, longitude, h3Index)
	if err != nil {
		s.logger.Error("Failed to update profile transaction: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update profile information",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully updated profile: ", update)
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    update,
	})
	return
}

func (s *HandlersServer) GetRolesFunc(ctx *gin.Context) {
	roles, err := s.store.GetRoles(ctx)
	if err != nil {
		s.logger.Error("Failed to get roles: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get roles",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    roles,
	})
}

func (s *HandlersServer) AddUserRoleFunc(ctx *gin.Context) {
	var req struct {
		RoleID int32 `json:"role_id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	roles, err := s.store.AddRole(ctx, authPayload.PublicID, req.RoleID)
	if err != nil {
		s.logger.Error("Failed to add role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to add role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    roles,
	})
}

func (s *HandlersServer) AddUserVerificationFunc(ctx *gin.Context) {
	var req struct {
		ProfileID        int64  `json:"profile_id" binding:"required,min=1"`
		OrganizationName string `json:"organization_name" binding:"required,min=2,max=200"`
		Email            string `json:"email" binding:"required,email"`
		PhoneNumber      string `json:"phone_number" binding:"required,min=10,max=15"`
		DocumentType     string `json:"document_type" binding:"required,min=2,max=50"`
		DocumentURL      string `json:"document_url" binding:"omitempty"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	//AddUserVerificationDocuments:
	organizerDetails, err := s.store.AddOrganizerVerificationDetails(ctx, req.ProfileID, req.OrganizationName, req.Email, req.PhoneNumber)
	if err != nil {
		s.logger.Error("Failed to Verify the organizer details: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to verify organizer details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	saveImageStruct := util.NewSaveImageStruct(s.logger)
	var emptyString string
	var documentVerification *models.Document
	if req.DocumentURL != emptyString {
		documentPath, err := saveImageStruct.SaveImageToFile([]byte(req.DocumentURL), "image")
		if err != nil {
			s.logger.Error("Failed to save avatar image: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to save document image",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		//Upload the documents:
		documentVerification, err = s.store.AddDocumentVerificationDetails(ctx, organizerDetails.ID, req.DocumentType, documentPath)
		if err != nil {
			s.logger.Error("Failed to verify document: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to verify document",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}
	}

	status := "pending"
	if organizerDetails.IsVerified {
		status = "verified"
	} else if organizerDetails.VerifiedAT != nil {
		status = "rejected"
	}

	var filePath string
	if documentVerification != nil {
		filePath = documentVerification.FilePath
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data": gin.H{
			"organizer_id":      organizerDetails.ID,
			"organization_name": organizerDetails.OrganizationName,
			"email":             organizerDetails.Email,
			"phone":             organizerDetails.PhoneNumber,
			"profile_id":        organizerDetails.ProfileID,
			"status":            status,
			"file_path":         filePath,
		},
	})

}

//Add the update status functionality

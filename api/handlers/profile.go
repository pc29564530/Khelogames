package handlers

import (
	"fmt"
	"khelogames/core/token"
	"khelogames/database/models"
	"khelogames/pkg"
	"khelogames/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) GetPlayerWithProfileFunc(ctx *gin.Context) {
	var req struct {
		PublicID string `json:"public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	playerProfile, err := s.store.GetPlayerWithProfile(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get player with profile", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get player with profile",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, playerProfile)
	return
}

func (s *HandlersServer) GetProfileByPublicIDFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get profile")
	var req struct {
		PublicID string `uri:"profile_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	profile, err := s.store.GetProfileByPublicID(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get profile",
		})
		return
	}
	s.logger.Debug("Successfully retrieved profile by profile public_id: ", profile)
	ctx.JSON(http.StatusOK, profile)
}

func (s *HandlersServer) GetProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get profile")
	var req struct {
		PublicID string `uri:"public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	profile, err := s.store.GetProfile(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get profile",
		})
		return
	}
	s.logger.Debug("Successfully retrieved profile: ", profile)
	ctx.JSON(http.StatusOK, profile)
}

func (s *HandlersServer) GetProfileByFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get profile")
	var req struct {
		PublicID string `uri:"public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	profile, err := s.store.GetProfile(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get profile",
		})
		return
	}
	s.logger.Info("Successfully retrieved profile: ", profile)
	ctx.JSON(http.StatusOK, profile)
}

type editProfileRequest struct {
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url,omitempty"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func (s *HandlersServer) UpdateProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update profile")
	var req editProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request data",
		})
		return
	}

	s.logger.Debug("Request JSON bind successful: ", req)

	latitude, err := strconv.ParseFloat(req.Latitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse latitude: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": fmt.Sprintf("Invalid latitude value: %v", req.Latitude),
		})
		return
	}
	longitude, err := strconv.ParseFloat(req.Longitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse longitude: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": fmt.Sprintf("Invalid longitude value: %v", req.Longitude),
		})
		return
	}
	var h3Index string
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	update, err := s.txStore.UpdateProfileTx(ctx, authPayload.PublicID, req.Bio, req.AvatarUrl, req.FullName, req.City, req.State, req.Country, latitude, longitude, h3Index)
	if err != nil {
		s.logger.Error("Failed to update profile transaction: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to update profile information",
		})
		return
	}

	s.logger.Info("Successfully updated profile: ", update)
	ctx.JSON(http.StatusAccepted, update)
	return
}

func (s *HandlersServer) GetRolesFunc(ctx *gin.Context) {
	roles, err := s.store.GetRoles(ctx)
	if err != nil {
		s.logger.Error("Failed to get roles: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get roles",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, roles)
}

func (s *HandlersServer) AddUserRoleFunc(ctx *gin.Context) {
	var req struct {
		RoleID int32 `json:"role_id"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request data",
		})
		return
	}
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	roles, err := s.store.AddRole(ctx, authPayload.PublicID, req.RoleID)
	if err != nil {
		s.logger.Error("Failed to get roles: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get roles",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, roles)
}

func (s *HandlersServer) AddUserVerificationFunc(ctx *gin.Context) {
	var req struct {
		ProfileID        int64  `json:"profile_id"`
		OrganizationName string `json:"organization_name"`
		Email            string `json:"email"`
		PhoneNumber      string `json:"phone_number"`
		DocumentType     string `json:"document_type"`
		DocumentURL      string `json:"document_url"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request data",
		})
		return
	}

	//AddUserVerificationDocuments:
	organizerDetails, err := s.store.AddOrganizerVerificationDetails(ctx, req.ProfileID, req.OrganizationName, req.Email, req.PhoneNumber)
	if err != nil {
		s.logger.Error("Failed to Verify the organizer details: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to verify organizer details",
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
				"code":    "INTERNAL_ERROR",
				"message": "Failed to save document image",
			})
			return
		}

		//Upload the documents:
		documentVerification, err = s.store.AddDocumentVerificationDetails(ctx, organizerDetails.ID, req.DocumentType, documentPath)
		if err != nil {
			s.logger.Error("Failed to verify document: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    "INTERNAL_ERROR",
				"message": "Failed to verify document",
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

	ctx.JSON(http.StatusAccepted, gin.H{
		"organizer_id":      organizerDetails.ID,
		"organization_name": organizerDetails.OrganizationName,
		"email":             organizerDetails.Email,
		"phone":             organizerDetails.PhoneNumber,
		"profile_id":        organizerDetails.ProfileID,
		"status":            status,
		"file_path":         documentVerification.FilePath,
	})

}

//Add the update status functionality

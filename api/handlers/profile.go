package handlers

import (
	"fmt"
	db "khelogames/database"
	"khelogames/database/models"
	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"

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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	playerProfile, err := s.store.GetPlayerWithProfile(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get player with profile", err)
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	profile, err := s.store.GetProfileByPublicID(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	fmt.Println("Profile; ", profile)
	s.logger.Info("Successfully retrieved profile by profile public_id: ", profile)
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	profile, err := s.store.GetProfile(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	fmt.Println("Profile; ", profile)
	s.logger.Info("Successfully retrieved profile: ", profile)
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	profile, err := s.store.GetProfile(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	fmt.Println("Profile; ", profile)
	s.logger.Info("Successfully retrieved profile: ", profile)
	ctx.JSON(http.StatusOK, profile)
}

type editProfileRequest struct {
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}

func (s *HandlersServer) UpdateProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update profile")
	var req editProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation: ", err)
		return
	}

	defer tx.Rollback()

	s.logger.Debug("Request JSON bind successful: ", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.EditProfileParams{
		PublicID:  authPayload.PublicID,
		Bio:       req.Bio,
		AvatarUrl: req.AvatarUrl,
	}

	updatedProfile, err := s.store.EditProfile(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update profile: ", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}

	_, err = s.store.UpdateUser(ctx, int32(updatedProfile.UserID), req.FullName)
	if err != nil {
		s.logger.Error("Failed to update the user full name: ", err)
		ctx.JSON(http.StatusNotAcceptable, err)
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	s.logger.Info("Successfully updated profile: ", updatedProfile)
	ctx.JSON(http.StatusAccepted, updatedProfile)
	return
}

func (s *HandlersServer) GetRolesFunc(ctx *gin.Context) {
	roles, err := s.store.GetRoles(ctx)
	if err != nil {
		s.logger.Error("Failed to get roles: ", err)
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
		return
	}
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	roles, err := s.store.AddRole(ctx, authPayload.PublicID, req.RoleID)
	if err != nil {
		s.logger.Error("Failed to get roles: ", err)
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
		return
	}

	//AddUserVerificationDocuments:
	organizerDetails, err := s.store.AddOrganizerVerificationDetails(ctx, req.ProfileID, req.OrganizationName, req.Email, req.PhoneNumber)
	if err != nil {
		s.logger.Error("Failed to Verify the organizer details: ", err)
		return
	}

	saveImageStruct := util.NewSaveImageStruct(s.logger)
	var emptyString string
	var documentVerification *models.Document
	if req.DocumentURL != emptyString {
		documentPath, err := saveImageStruct.SaveImageToFile([]byte(req.DocumentURL), "image")
		if err != nil {
			s.logger.Error("Failed to save avatar image: ", err)
			return
		}

		//Upload the documents:
		documentVerification, err = s.store.AddDocumentVerificationDetails(ctx, organizerDetails.ID, req.DocumentType, documentPath)
		if err != nil {
			s.logger.Error("Failed to verify document: ", err)
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

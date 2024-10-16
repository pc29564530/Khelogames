package handlers

import (
	"encoding/base64"
	db "khelogames/database"
	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type createProfileRequest struct {
	FullName  string `json:"full_name,omitempty"`
	Bio       string `json:"bio,omitempty"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}

func (s *HandlersServer) CreateProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create profile")
	var req createProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: ", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateProfileParams{
		Owner:     authPayload.Username,
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarUrl: req.AvatarUrl,
	}

	profile, err := s.store.CreateProfile(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create profile: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully created profile: ", profile)
	ctx.JSON(http.StatusOK, profile)
	return
}

type getProfileRequest struct {
	Owner string `uri:"owner"`
}

func (s *HandlersServer) GetProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get profile")
	var req getProfileRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	profile, err := s.store.GetProfile(ctx, req.Owner)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
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

	b64data := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		s.logger.Error("Failed to decode avatar string: ", err)
		return
	}
	s.logger.Debug("Avatar string decoded successfully")

	saveImageStruct := util.NewSaveImageStruct(s.logger)
	var avatarPath string
	mediaType := "image"
	if req.AvatarUrl != "" {
		avatarPath, err = saveImageStruct.SaveImageToFile(avatarData, mediaType)
		if err != nil {
			s.logger.Error("Failed to save avatar image: ", err)
			return
		}
		s.logger.Debug("Avatar saved successfully at ", avatarPath)
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	profile, err := s.store.GetProfile(ctx, authPayload.Username)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.EditProfileParams{
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarUrl: avatarPath,
		ID:        profile.ID,
	}

	updatedProfile, err := s.store.EditProfile(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update profile: ", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
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

type editFullNameRequest struct {
	FullName string `json:"full_name"`
}

func (s *HandlersServer) UpdateFullNameFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update full name")
	var req editFullNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: ", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	profileFullName, err := s.store.UpdateFullName(ctx, authPayload.Username, req.FullName)
	if err != nil {
		s.logger.Error("Failed to update full name: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully updated full name: ", profileFullName)
	ctx.JSON(http.StatusAccepted, profileFullName)
	return
}

type editBioRequest struct {
	Bio string `json:"bio"`
}

func (s *HandlersServer) UpdateBioFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update bio")
	var req editBioRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: ", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.UpdateBioParams{
		Owner: authPayload.Username,
		Bio:   req.Bio,
	}

	profileBio, err := s.store.UpdateBio(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update bio: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully updated bio: ", profileBio)
	ctx.JSON(http.StatusAccepted, profileBio)
	return
}

type editAvatarUrlRequest struct {
	AvatarUrl string `json:"avatar_url,omitempty"`
}

func (s *HandlersServer) UpdateAvatarUrlFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update avatar URL")
	var req editAvatarUrlRequest
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

	b64data := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		s.logger.Error("Failed to decode avatar string: ", err)
		return
	}
	s.logger.Debug("Avatar string decoded successfully")
	saveImageStruct := util.NewSaveImageStruct(s.logger)
	mediaType := "image"
	path, err := saveImageStruct.SaveImageToFile(avatarData, mediaType)
	if err != nil {
		s.logger.Error("Failed to save avatar image: ", err)
		return
	}
	s.logger.Debug("Avatar saved successfully at ", path)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.UpdateAvatarParams{
		Owner:     authPayload.Username,
		AvatarUrl: path,
	}

	profileAvatarUrl, err := s.store.UpdateAvatar(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update avatar URL: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	s.logger.Info("Successfully updated avatar URL: ", profileAvatarUrl)
	ctx.JSON(http.StatusAccepted, profileAvatarUrl)
	return
}

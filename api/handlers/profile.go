package handlers

import (
	"encoding/base64"
	db "khelogames/db/sqlc"
	"khelogames/logger"
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
	CoverUrl  string `json:"cover_url,omitempty"`
}

type ProfileServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewProfileServer(store *db.Store, logger *logger.Logger) *ProfileServer {
	return &ProfileServer{store: store, logger: logger}
}

func (s *ProfileServer) CreateProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create profile")
	var req createProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: %v", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateProfileParams{
		Owner:     authPayload.Username,
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarUrl: req.AvatarUrl,
		CoverUrl:  req.CoverUrl,
	}

	profile, err := s.store.CreateProfile(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create profile: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully created profile: %v", profile)
	ctx.JSON(http.StatusOK, profile)
	return
}

type getProfileRequest struct {
	Owner string `uri:"owner"`
}

func (s *ProfileServer) GetProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get profile")
	var req getProfileRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind URI: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request URI bind successful: %v", req)

	profile, err := s.store.GetProfile(ctx, req.Owner)
	if err != nil {
		s.logger.Error("Failed to get profile: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Info("Successfully retrieved profile: %v", profile)
	ctx.JSON(http.StatusOK, profile)
	return
}

type editProfileRequest struct {
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url,omitempty"`
	CoverUrl  string `json:"cover_url,omitempty"`
}

func (s *ProfileServer) UpdateProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update profile")
	var req editProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: %v", req)

	b64data := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		s.logger.Error("Failed to decode avatar string: %v", err)
		return
	}
	s.logger.Debug("Avatar string decoded successfully")

	b64data = req.CoverUrl[strings.IndexByte(req.CoverUrl, ',')+1:]

	coverData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		s.logger.Error("Failed to decode cover string: %v", err)
		return
	}
	s.logger.Debug("Cover string decoded successfully")
	saveImageStruct := util.NewSaveImageStruct(s.logger)
	var avatarPath string
	mediaType := "image"
	if req.AvatarUrl != "" {
		avatarPath, err = saveImageStruct.SaveImageToFile(avatarData, mediaType)
		if err != nil {
			s.logger.Error("Failed to save avatar image: %v", err)
			return
		}
		s.logger.Debug("Avatar saved successfully at %s", avatarPath)
	}

	var coverPath string
	if req.CoverUrl != "" {
		coverPath, err = saveImageStruct.SaveImageToFile(coverData, mediaType)
		if err != nil {
			s.logger.Error("Failed to save cover image: %v", err)
			return
		}
		s.logger.Debug("Cover saved successfully at %s", coverPath)
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	profile, err := s.store.GetProfile(ctx, authPayload.Username)
	if err != nil {
		s.logger.Error("Failed to get profile: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.EditProfileParams{
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarUrl: avatarPath,
		CoverUrl:  coverPath,
		ID:        profile.ID,
	}

	updatedProfile, err := s.store.EditProfile(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update profile: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	s.logger.Info("Successfully updated profile: %v", updatedProfile)
	ctx.JSON(http.StatusAccepted, updatedProfile)
	return
}

type editFullNameRequest struct {
	FullName string `json:"full_name"`
}

func (s *ProfileServer) UpdateFullNameFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update full name")
	var req editFullNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: %v", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.UpdateFullNameParams{
		Owner:    authPayload.Username,
		FullName: req.FullName,
	}

	profileFullName, err := s.store.UpdateFullName(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update full name: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully updated full name: %v", profileFullName)
	ctx.JSON(http.StatusAccepted, profileFullName)
	return
}

type editBioRequest struct {
	Bio string `json:"bio"`
}

func (s *ProfileServer) UpdateBioFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update bio")
	var req editBioRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: %v", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.UpdateBioParams{
		Owner: authPayload.Username,
		Bio:   req.Bio,
	}

	profileBio, err := s.store.UpdateBio(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update bio: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully updated bio: %v", profileBio)
	ctx.JSON(http.StatusAccepted, profileBio)
	return
}

type editAvatarUrlRequest struct {
	AvatarUrl string `json:"avatar_url,omitempty"`
}

func (s *ProfileServer) UpdateAvatarUrlFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update avatar URL")
	var req editAvatarUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: %v", req)

	b64data := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		s.logger.Error("Failed to decode avatar string: %v", err)
		return
	}
	s.logger.Debug("Avatar string decoded successfully")
	saveImageStruct := util.NewSaveImageStruct(s.logger)
	mediaType := "image"
	path, err := saveImageStruct.SaveImageToFile(avatarData, mediaType)
	if err != nil {
		s.logger.Error("Failed to save avatar image: %v", err)
		return
	}
	s.logger.Debug("Avatar saved successfully at %s", path)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.UpdateAvatarParams{
		Owner:     authPayload.Username,
		AvatarUrl: path,
	}

	profileAvatarUrl, err := s.store.UpdateAvatar(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update avatar URL: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully updated avatar URL: %v", profileAvatarUrl)
	ctx.JSON(http.StatusAccepted, profileAvatarUrl)
	return
}

type editCoverUrlRequest struct {
	CoverUrl string `json:"cover_url,omitempty"`
}

func (s *ProfileServer) UpdateCoverUrlFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update cover URL")
	var req editCoverUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: %v", req)

	b64data := req.CoverUrl[strings.IndexByte(req.CoverUrl, ',')+1:]

	coverData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		s.logger.Error("Failed to decode cover string: %v", err)
		return
	}
	s.logger.Debug("Cover string decoded successfully")
	saveImageStruct := util.NewSaveImageStruct(s.logger)
	mediaType := "image"
	path, err := saveImageStruct.SaveImageToFile(coverData, mediaType)
	if err != nil {
		s.logger.Error("Failed to save cover image: %v", err)
		return
	}
	s.logger.Debug("Cover saved successfully at %s", path)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.UpdateCoverParams{
		Owner:    authPayload.Username,
		CoverUrl: path,
	}

	profileCoverUrl, err := s.store.UpdateCover(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update cover URL: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully updated cover URL: %v", profileCoverUrl)
	ctx.JSON(http.StatusAccepted, profileCoverUrl)
	return
}

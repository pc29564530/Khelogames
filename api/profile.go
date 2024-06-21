package api

import (
	"encoding/base64"
	db "khelogames/db/sqlc"
	"khelogames/token"
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

func (server *Server) createProfile(ctx *gin.Context) {
	var req createProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateProfileParams{
		Owner:     authPayload.Username,
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarUrl: req.AvatarUrl,
		CoverUrl:  req.CoverUrl,
	}

	profile, err := server.store.CreateProfile(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to create profile: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	server.logger.Info("Successfully created profile")
	ctx.JSON(http.StatusOK, profile)
	return
}

type getProfileRequest struct {
	Owner string `uri:"owner"`
}

func (server *Server) getProfile(ctx *gin.Context) {
	var req getProfileRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	profile, err := server.store.GetProfile(ctx, req.Owner)
	if err != nil {
		server.logger.Error("Failed to get profile: %v", err)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	server.logger.Info("Successfully created profile")
	ctx.JSON(http.StatusOK, profile)
	return
}

type editProfileRequest struct {
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url,omitempty"`
	CoverUrl  string `json:"cover_url,omitempty"`
}

func (server *Server) updateProfile(ctx *gin.Context) {

	var req editProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	b64data := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		server.logger.Error("Failed to decode avatar string: %v", err)
		return
	}

	b64data = req.CoverUrl[strings.IndexByte(req.CoverUrl, ',')+1:]

	coverData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		server.logger.Error("Failed to decode cover string: %v", err)
		return
	}

	var avatarPath string
	mediaType := "image"
	if req.AvatarUrl != "" {
		avatarPath, err = saveImageToFile(avatarData, mediaType)
		if err != nil {
			server.logger.Error("Failed to create the avatar path: %v", err)
			return
		}
	}
	var coverPath string
	if req.CoverUrl != "" {
		coverPath, err = saveImageToFile(coverData, mediaType)
		if err != nil {
			server.logger.Error("Failed to create cover path: %v", err)
			return
		}
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	profile, err := server.store.GetProfile(ctx, authPayload.Username)
	if err != nil {
		server.logger.Error("Failed to get profile: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.EditProfileParams{
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarUrl: avatarPath,
		CoverUrl:  coverPath,
		ID:        profile.ID,
	}

	updatedProfile, err := server.store.EditProfile(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to edit profile: %v", err)
		ctx.JSON(http.StatusNotAcceptable, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, updatedProfile)
	return
}

type editFullNameRequest struct {
	FullName string `json:"full_name"`
}

func (server *Server) updateFullName(ctx *gin.Context) {
	var req editFullNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UpdateFullNameParams{
		Owner:    authPayload.Username,
		FullName: req.FullName,
	}

	profileFullName, err := server.store.UpdateFullName(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to update full name: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	server.logger.Info("Successfully updated full name")
	ctx.JSON(http.StatusAccepted, profileFullName)
	return
}

type editBioRequest struct {
	Bio string `json:"bio"`
}

func (server *Server) updateBio(ctx *gin.Context) {
	var req editBioRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UpdateBioParams{
		Owner: authPayload.Username,
		Bio:   req.Bio,
	}

	profileBio, err := server.store.UpdateBio(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to update bio: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, profileBio)
	return
}

type editAvatarUrlRequest struct {
	AvatarUrl string `json:"avatar_url,omitempty"`
}

func (server *Server) updateAvatarUrl(ctx *gin.Context) {
	var req editAvatarUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to update avatar url: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	b64data := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		server.logger.Error("Failed to decode avatar: %v", err)
		return
	}
	mediaType := "image"
	path, err := saveImageToFile(avatarData, mediaType)
	if err != nil {
		server.logger.Error("Failed to create avatar file: %v", err)
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UpdateAvatarParams{
		Owner:     authPayload.Username,
		AvatarUrl: path,
	}

	profileAvatarUrl, err := server.store.UpdateAvatar(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to update avatar: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, profileAvatarUrl)
	return
}

type editCoverUrlRequest struct {
	CoverUrl string `json:"cover_url,omitempty"`
}

func (server *Server) updateCoverUrl(ctx *gin.Context) {
	var req editCoverUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	b64data := req.CoverUrl[strings.IndexByte(req.CoverUrl, ',')+1:]

	coverData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		server.logger.Error("Failed to decode covert url: %v", err)
		return
	}
	mediaType := "image"
	path, err := saveImageToFile(coverData, mediaType)
	if err != nil {
		server.logger.Error("Failed to create file: %v", err)
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UpdateCoverParams{
		Owner:    authPayload.Username,
		CoverUrl: path,
	}

	profileCoverUrl, err := server.store.UpdateCover(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to update cover: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, profileCoverUrl)
	return
}

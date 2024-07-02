package handlers

import (
	"encoding/base64"
	"fmt"
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
	var req createProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

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
		fmt.Errorf("Failed to create profile: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println("Successfully created profile")
	ctx.JSON(http.StatusOK, profile)
	return
}

type getProfileRequest struct {
	Owner string `uri:"owner"`
}

func (s *ProfileServer) GetProfileFunc(ctx *gin.Context) {
	fmt.Println("Line no 24 Profile")
	var req getProfileRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	profile, err := s.store.GetProfile(ctx, req.Owner)
	if err != nil {
		fmt.Errorf("Failed to get profile: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	fmt.Println("Successfully created profile")
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

	var req editProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	b64data := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Errorf("Failed to decode avatar string: %v", err)
		return
	}

	b64data = req.CoverUrl[strings.IndexByte(req.CoverUrl, ',')+1:]

	coverData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Errorf("Failed to decode cover string: %v", err)
		return
	}

	var avatarPath string
	mediaType := "image"
	if req.AvatarUrl != "" {
		avatarPath, err = util.SaveImageToFile(avatarData, mediaType)
		if err != nil {
			fmt.Errorf("Failed to create the avatar path: %v", err)
			return
		}
	}
	var coverPath string
	if req.CoverUrl != "" {
		coverPath, err = util.SaveImageToFile(coverData, mediaType)
		if err != nil {
			fmt.Errorf("Failed to create cover path: %v", err)
			return
		}
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	profile, err := s.store.GetProfile(ctx, authPayload.Username)
	if err != nil {
		fmt.Errorf("Failed to get profile: %v", err)
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
		fmt.Errorf("Failed to edit profile: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}

	ctx.JSON(http.StatusAccepted, updatedProfile)
	return
}

type editFullNameRequest struct {
	FullName string `json:"full_name"`
}

func (s *ProfileServer) UpdateFullNameFunc(ctx *gin.Context) {
	var req editFullNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.UpdateFullNameParams{
		Owner:    authPayload.Username,
		FullName: req.FullName,
	}

	profileFullName, err := s.store.UpdateFullName(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update full name: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println("Successfully updated full name")
	ctx.JSON(http.StatusAccepted, profileFullName)
	return
}

type editBioRequest struct {
	Bio string `json:"bio"`
}

func (s *ProfileServer) UpdateBioFunc(ctx *gin.Context) {
	var req editBioRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.UpdateBioParams{
		Owner: authPayload.Username,
		Bio:   req.Bio,
	}

	profileBio, err := s.store.UpdateBio(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update bio: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusAccepted, profileBio)
	return
}

type editAvatarUrlRequest struct {
	AvatarUrl string `json:"avatar_url,omitempty"`
}

func (s *ProfileServer) UpdateAvatarUrlFunc(ctx *gin.Context) {
	var req editAvatarUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to update avatar url: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	b64data := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Errorf("Failed to decode avatar: %v", err)
		return
	}
	mediaType := "image"
	path, err := util.SaveImageToFile(avatarData, mediaType)
	if err != nil {
		fmt.Errorf("Failed to create avatar file: %v", err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.UpdateAvatarParams{
		Owner:     authPayload.Username,
		AvatarUrl: path,
	}

	profileAvatarUrl, err := s.store.UpdateAvatar(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update avatar: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusAccepted, profileAvatarUrl)
	return
}

type editCoverUrlRequest struct {
	CoverUrl string `json:"cover_url,omitempty"`
}

func (s *ProfileServer) UpdateCoverUrlFunc(ctx *gin.Context) {
	var req editCoverUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	b64data := req.CoverUrl[strings.IndexByte(req.CoverUrl, ',')+1:]

	coverData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Errorf("Failed to decode covert url: %v", err)
		return
	}
	mediaType := "image"
	path, err := util.SaveImageToFile(coverData, mediaType)
	if err != nil {
		fmt.Errorf("Failed to create file: %v", err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.UpdateCoverParams{
		Owner:    authPayload.Username,
		CoverUrl: path,
	}

	profileCoverUrl, err := s.store.UpdateCover(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update cover: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusAccepted, profileCoverUrl)
	return
}

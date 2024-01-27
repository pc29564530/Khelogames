package api

import (
	"encoding/base64"
	"fmt"
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	profile, err := server.store.GetProfile(ctx, req.Owner)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	b64data := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Println("unable to decode avatar :", err)
		return
	}

	b64data = req.CoverUrl[strings.IndexByte(req.CoverUrl, ',')+1:]

	coverData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Println("unable to decode cover  :", err)
		return
	}

	var avatarPath string
	mediaType := "image"
	if req.AvatarUrl != "" {
		avatarPath, err = saveImageToFile(avatarData, mediaType)
		if err != nil {
			fmt.Println("unable to create a avatar file")
			return
		}
	}
	var coverPath string
	if req.CoverUrl != "" {
		coverPath, err = saveImageToFile(coverData, mediaType)
		if err != nil {
			fmt.Println("unable to create a cover file")
			return
		}
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	profile, err := server.store.GetProfile(ctx, authPayload.Username)
	if err != nil {
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	b64data := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Println("unable to decode avatar :", err)
		return
	}
	mediaType := "image"
	path, err := saveImageToFile(avatarData, mediaType)
	if err != nil {
		fmt.Println("unable to create a avatar file")
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UpdateAvatarParams{
		Owner:     authPayload.Username,
		AvatarUrl: path,
	}

	profileAvatarUrl, err := server.store.UpdateAvatar(ctx, arg)
	if err != nil {
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	b64data := req.CoverUrl[strings.IndexByte(req.CoverUrl, ',')+1:]

	coverData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Println("unable to decode cover  :", err)
		return
	}
	mediaType := "image"
	path, err := saveImageToFile(coverData, mediaType)
	if err != nil {
		fmt.Println("unable to create a cover file")
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UpdateCoverParams{
		Owner:    authPayload.Username,
		CoverUrl: path,
	}

	profileCoverUrl, err := server.store.UpdateCover(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, profileCoverUrl)
	return
}

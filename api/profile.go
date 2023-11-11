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
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url"`
	CoverUrl  string `json:"cover_url"`
}

func (server *Server) createProfile(ctx *gin.Context) {
	var req createProfileRequest
	fmt.Println("Hello")
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println("FullName:")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	fmt.Println("Username: ", authPayload.Username)
	arg := db.CreateProfileParams{
		Owner:     authPayload.Username,
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarUrl: req.AvatarUrl,
		CoverUrl:  req.CoverUrl,
	}

	fmt.Println("Arg: ", arg)

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
	fmt.Println("Hello")
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println("Profile Owner: ", req.Owner)

	profile, err := server.store.GetProfile(ctx, req.Owner)
	fmt.Println("Profile: ", profile)
	fmt.Println("Error: ", err)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, profile)
	return

}

type editProfileRequest struct {
	FullName  string `json:"full_name,omitempty"`
	Bio       string `json:"bio,omitempty"`
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

	b64Avatar := req.AvatarUrl[strings.IndexByte(req.AvatarUrl, ',')+1:]

	avatarData, err := base64.StdEncoding.DecodeString(b64Avatar)
	if err != nil {
		fmt.Println("uanble to decode :", err)
		return
	}

	b64Cover := req.CoverUrl[strings.IndexByte(req.CoverUrl, ',')+1:]

	coverData, err := base64.StdEncoding.DecodeString(b64Cover)
	if err != nil {
		fmt.Println("uanble to decode :", err)
		return
	}

	avatarPath, err := saveImageToFile(avatarData)
	if err != nil {
		fmt.Println("unable to create a file")
		return
	}

	coverPath, err := saveImageToFile(coverData)
	if err != nil {
		fmt.Println("unable to create a file")
		return
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

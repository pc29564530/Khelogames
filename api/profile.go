package api

import (
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
)

type createProfileRequest struct {
	FullName       string `json:"full_name"`
	Bio            string `json:"bio,omitempty"`
	FollowingOwner int64  `json:"following_owner"`
	FollowerOwner  int64  `json:"follower_owner"`
	AvatarUrl      string `json:"avatar_url"`
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
		Owner:          authPayload.Username,
		FullName:       req.FullName,
		Bio:            req.Bio,
		FollowingOwner: req.FollowingOwner,
		FollowerOwner:  req.FollowerOwner,
		AvatarUrl:      req.AvatarUrl,
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

type updateBioRequest struct {
	Bio string `json:"bio"`
	ID  int64  `json:"id"`
}

func (server *Server) updateBio(ctx *gin.Context) {
	var req updateBioRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateProfileBioParams{
		Bio: req.Bio,
		ID:  req.ID,
	}

	bio, err := server.store.UpdateProfileBio(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bio)
	return
}

type updateAvatarUrlRequest struct {
	AvatarUrl string `json:"avatar_url"`
	ID        int64  `json:"id"`
}

func (server *Server) updateAvatarUrl(ctx *gin.Context) {
	var req updateAvatarUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateProfileAvatarParams{
		AvatarUrl: req.AvatarUrl,
		ID:        req.ID,
	}

	avatarUrl, err := server.store.UpdateProfileAvatar(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, avatarUrl)
	return
}

type updateFullNameRequest struct {
	FullName string `json:"full_name"`
	ID       int64  `json:"id"`
}

func (server *Server) updateFullName(ctx *gin.Context) {
	var req updateFullNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateProfileFullNameParams{
		FullName: req.FullName,
		ID:       req.ID,
	}

	bio, err := server.store.UpdateProfileFullName(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bio)
	return
}

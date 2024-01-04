package api

import (
	"database/sql"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createCommunitiesRequest struct {
	Owner         string `json:"owner"`
	CommunityName string `json:"communityName"`
	Description   string `json:"description"`
	CommunityType string `json:"communityType"`
}

// Create communities function
func (server *Server) createCommunites(ctx *gin.Context) {
	var req createCommunitiesRequest
	fmt.Println(req)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateCommunityParams{
		Owner:           authPayload.Username,
		CommunitiesName: req.CommunityName,
		Description:     req.Description,
		CommunityType:   req.CommunityType,
	}

	communities, err := server.store.CreateCommunity(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, communities)
	return
}

type getCommunityRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type getCommunityResponse struct {
	CommunitiesName string    `json:"communityName"`
	Description     string    `json:"description"`
	CommunityType   string    `json:"communityType"`
	CreatedAt       time.Time `json:"created_at"`
}

// get Community by id.
func (server *Server) getCommunity(ctx *gin.Context) {
	var req getCommunityRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	community, err := server.store.GetCommunity(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := getCommunityResponse{
		CommunitiesName: community.CommunitiesName,
		Description:     community.Description,
		CommunityType:   community.CommunityType,
		CreatedAt:       community.CreatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
	return
}

// Get all communities by owner.
func (server *Server) getAllCommunities(ctx *gin.Context) {

	user, err := server.store.GetAllCommunities(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
	return
}

// Get all users that have joined a particular communities
type getCommunitiesMemberRequest struct {
	CommunitiesName string `uri:"communities_name"`
}

func (server *Server) getCommunitiesMember(ctx *gin.Context) {
	var req getCommunitiesMemberRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	usersList, err := server.store.GetCommunitiesMember(ctx, req.CommunitiesName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, usersList)
}

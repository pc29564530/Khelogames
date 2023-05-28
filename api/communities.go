package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"net/http"
	"time"
)

type createCommunitiesRequest struct {
	CommunitiesName string `json:"communities_name"`
	Description     string `json:"description"`
	CommunityType   string `json:"community_type"`
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

	arg := db.CreateCommunityParams{
		CommunitiesName: req.CommunitiesName,
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
	CommunitiesName string    `json:"communities_name"`
	Description     string    `json:"description"`
	CommunityType   string    `json:"community_type"`
	CreatedAt       time.Time `json:"created_at"`
}

// Open selected community
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

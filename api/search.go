package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type searchByCommunityNameRequest struct {
	CommunityName string `json:"communities_name"`
}

func (server *Server) searchByCommunityName(ctx *gin.Context) {
	var req searchByCommunityNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	communityName, err := server.store.SearchByCommunityName(ctx, req.CommunityName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, communityName)
	return
}

type searchByFullNameRequest struct {
	FullName string `json:"full_name"`
}

func (server *Server) searchByFullName(ctx *gin.Context) {
	var req searchByFullNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	profile, err := server.store.SearchByFullName(ctx, req.FullName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, profile)
	return
}

type searchCommunityByCommunityTypeRequest struct {
	CommunityType string `json:"community_type"`
}

func (server *Server) searchCommunityByCommunityType(ctx *gin.Context) {
	var req searchCommunityByCommunityTypeRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	profile, err := server.store.SearchCommunityByCommunityType(ctx, req.CommunityType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, profile)
	return
}

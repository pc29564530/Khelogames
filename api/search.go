package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type searchByCommunityNameRequest struct {
	CommunityName string `uri:"communities_name"`
}

func (server *Server) searchByCommunityName(ctx *gin.Context) {
	var req searchByCommunityNameRequest
	err := ctx.ShouldBindUri(&req)
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
	FullName string `uri:"full_name"`
}

func (server *Server) searchByFullName(ctx *gin.Context) {
	var req searchByFullNameRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	profile, err := server.store.SearchByCommunityName(ctx, req.FullName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, profile)
	return
}

type searchCommunityByCommunityTypeRequest struct {
	CommunityType string `uri:"community_type"`
}

func (server *Server) searchCommunityByCommunityType(ctx *gin.Context) {
	var req searchCommunityByCommunityTypeRequest
	err := ctx.ShouldBindUri(&req)
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

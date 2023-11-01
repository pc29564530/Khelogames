package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type searchByCommunityNameRequest struct {
	CommunityName string `json:"communities_name"`
}

func (server *Server) searchByCommunityName(ctx *gin.Context) {
	fmt.Println("Line no 14")
	var req searchByCommunityNameRequest
	fmt.Println("Line no 15")
	fmt.Println(req.CommunityName)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println("Line no 21")
	fmt.Println("CommunityName: %s", req.CommunityName)

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
	fmt.Println("Line no 63")

	var req searchCommunityByCommunityTypeRequest
	fmt.Println("Line no 66")
	fmt.Println("CommunityType", req)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println("Line no 70")
	fmt.Println("communityType: ", req.CommunityType)
	profile, err := server.store.SearchCommunityByCommunityType(ctx, req.CommunityType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	fmt.Println("Line no 77")
	fmt.Println("Profile: ", profile)
	ctx.JSON(http.StatusOK, profile)
	return
}

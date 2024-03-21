package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addClubMemberRequest struct {
	ClubName   string `json:"club_name"`
	ClubMember string `json:"club_member"`
}

func (server *Server) addClubMember(ctx *gin.Context) {
	var req addClubMemberRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.AddClubMemberParams{
		ClubName:   req.ClubName,
		ClubMember: req.ClubMember,
	}

	members, err := server.store.AddClubMember(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, members)
	return
}

type getClubMemberRequest struct {
	ClubName string `uri:"club_name"`
}

func (server *Server) getClubMember(ctx *gin.Context) {
	var req getClubMemberRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println("Club Name: ", req.ClubName)

	members, err := server.store.GetClubMember(ctx, req.ClubName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, members)
	return
}

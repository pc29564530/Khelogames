package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addClubMemberRequest struct {
	ClubID   int64 `json:"club_id"`
	PlayerID int64 `json:"player_id"`
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
		ClubID:   req.ClubID,
		PlayerID: req.PlayerID,
	}

	fmt.Println("Arg: lone no 29: ", arg)

	members, err := server.store.AddClubMember(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, members)
	return
}

type getClubMemberRequest struct {
	ClubID int64 `json:"club_id"`
}

func (server *Server) getClubMember(ctx *gin.Context) {
	clubIDStr := ctx.Query("club_id")
	clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
	if err != nil {
		fmt.Println("unable to parse the clubIDStr: ", err)
		return
	}
	fmt.Println("Line no 53")
	members, err := server.store.GetClubMember(ctx, clubID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, members)
	return
}

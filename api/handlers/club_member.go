package handlers

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClubMemberServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewClubMemberServer(store *db.Store, logger *logger.Logger) *ClubMemberServer {
	return &ClubMemberServer{store: store, logger: logger}
}

type addClubMemberRequest struct {
	ClubID   int64 `json:"club_id"`
	PlayerID int64 `json:"player_id"`
}

func (s *ClubMemberServer) AddClubMemberFunc(ctx *gin.Context) {
	var req addClubMemberRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.AddClubMemberParams{
		ClubID:   req.ClubID,
		PlayerID: req.PlayerID,
	}

	members, err := s.store.AddClubMember(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to add club member: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, members)
	return
}

type getClubMemberRequest struct {
	ClubID int64 `json:"club_id"`
}

func (s *ClubMemberServer) GetClubMemberFunc(ctx *gin.Context) {
	clubIDStr := ctx.Query("club_id")
	clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse club id string: %v", err)
		return
	}

	members, err := s.store.GetClubMember(ctx, clubID)
	if err != nil {
		fmt.Errorf("Failed to get club member: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, members)
	return
}

package clubs

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addClubMemberRequest struct {
	ClubID   int64 `json:"club_id"`
	PlayerID int64 `json:"player_id"`
}

func (s *ClubServer) AddClubMemberFunc(ctx *gin.Context) {
	var req addClubMemberRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
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
		s.logger.Error("Failed to add club member: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Info("successfully added member to the club")
	ctx.JSON(http.StatusAccepted, members)
	return
}

type getClubMemberRequest struct {
	ClubID int64 `json:"club_id"`
}

func (s *ClubServer) GetClubMemberFunc(ctx *gin.Context) {
	clubIDStr := ctx.Query("club_id")
	clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse club id string: %v", err)
		return
	}
	s.logger.Debug("get club id from reqeust: %v", clubID)

	members, err := s.store.GetClubMember(ctx, clubID)
	if err != nil {
		s.logger.Error("Failed to get club member: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	s.logger.Info("successfully get club member")

	ctx.JSON(http.StatusAccepted, members)
	return
}

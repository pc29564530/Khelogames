package clubs

import (
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateAvatarUrlRequest struct {
	AvatarUrl string `json:"avatar_url"`
	ClubName  string `json:"club_name"`
}

func (s *ClubServer) UpdateClubAvatarUrlFunc(ctx *gin.Context) {
	var req updateAvatarUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateAvatarUrlParams{
		AvatarUrl: req.AvatarUrl,
		ClubName:  req.ClubName,
	}

	response, err := s.store.UpdateAvatarUrl(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update avatar url: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateClubNameRequest struct {
	AvatarUrl string `json:"avatar_url"`
	ClubName  string `json:"club_name"`
}

func (s *ClubServer) updateClubName(ctx *gin.Context) {
	var req updateClubNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	arg := db.UpdateClubNameParams{
		ClubName: req.ClubName,
	}

	response, err := s.store.UpdateClubName(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateClubSport struct {
	ClubName string `json:"club_name"`
	Sport    string `json:"sport"`
}

func (s *ClubServer) UpdateClubSportFunc(ctx *gin.Context) {
	var req updateClubSport
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	arg := db.UpdateClubSportParams{
		Sport:    req.Sport,
		ClubName: req.ClubName,
	}

	response, err := s.store.UpdateClubSport(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update club sport: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

package handlers

import (
	"database/sql"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CommunityServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCommunityServer(store *db.Store, logger *logger.Logger) *CommunityServer {
	return &CommunityServer{store: store, logger: logger}
}

type createCommunitiesRequest struct {
	CommunityName string `json:"communityName"`
	Description   string `json:"description"`
	CommunityType string `json:"communityType"`
}

// Create communities function
func (s *CommunityServer) CreateCommunitesFunc(ctx *gin.Context) {
	var req createCommunitiesRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateCommunityParams{
		Owner:           authPayload.Username,
		CommunitiesName: req.CommunityName,
		Description:     req.Description,
		CommunityType:   req.CommunityType,
	}

	communities, err := s.store.CreateCommunity(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create community: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
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
func (s *CommunityServer) GetCommunityFunc(ctx *gin.Context) {
	var req getCommunityRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	community, err := s.store.GetCommunity(ctx, req.ID)
	if err != nil {
		fmt.Errorf("Failed to get community: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
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
func (s *CommunityServer) GetAllCommunitiesFunc(ctx *gin.Context) {

	user, err := s.store.GetAllCommunities(ctx)
	if err != nil {
		fmt.Errorf("Failed to  get communities: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	ctx.JSON(http.StatusOK, user)
	return
}

// Get all users that have joined a particular communities
type getCommunitiesMemberRequest struct {
	CommunitiesName string `uri:"communities_name"`
}

func (s *CommunityServer) GetCommunitiesMemberFunc(ctx *gin.Context) {
	var req getCommunitiesMemberRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	usersList, err := s.store.GetCommunitiesMember(ctx, req.CommunitiesName)
	if err != nil {
		fmt.Errorf("Failed to get community member: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	ctx.JSON(http.StatusOK, usersList)
}

type getCommunityByCommunityNameRequest struct {
	CommunitiesName string `uri:"communities_name"`
}

func (s *CommunityServer) GetCommunityByCommunityNameFunc(ctx *gin.Context) {
	var req getCommunityByCommunityNameRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row error %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	usersList, err := s.store.GetCommunityByCommunityName(ctx, req.CommunitiesName)
	if err != nil {
		fmt.Errorf("Failed to get community by community name: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, usersList)
}

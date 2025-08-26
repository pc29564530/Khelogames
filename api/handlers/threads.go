package handlers

import (
	db "khelogames/database"
	"khelogames/pkg"
	"khelogames/token"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createThreadRequest struct {
	CommunityPublicID string `json:"community_public_id,omitempty"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	MediaType         string `json:"mediaType,omitempty"`
	MediaURL          string `json:"mediaURL,omitempty"`
}

func (s *HandlersServer) CreateThreadFunc(ctx *gin.Context) {
	var req createThreadRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transaction: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	var communityPublicID *uuid.UUID
	if req.CommunityPublicID != "" && req.CommunityPublicID != "null" {
		parsed, err := uuid.Parse(req.CommunityPublicID)
		if err != nil {
			s.logger.Error("Failed to parse community public id: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid community ID"})
			return
		}
		communityPublicID = &parsed
	}

	//function for uploading a image or video
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateThreadParams{
		UserPublicID:      authPayload.PublicID,
		CommunityPublicID: communityPublicID,
		Title:             req.Title,
		Content:           req.Content,
		MediaType:         req.MediaType,
		MediaUrl:          req.MediaURL,
	}

	thread, err := s.store.CreateThread(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create new thread ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create thread"})
		return
	}

	users, err := s.store.GetProfileByUserID(ctx, thread.UserID)
	if err != nil {
		s.logger.Error("Failed to get user profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit the transaction: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	threadResponse := map[string]interface{}{
		"id":            thread.ID,
		"public_id":     thread.PublicID,
		"user_id":       thread.UserID,
		"community_id":  thread.CommunityID,
		"title":         thread.Title,
		"content":       thread.Content,
		"media_url":     thread.MediaUrl,
		"media_type":    thread.MediaType,
		"like_count":    thread.LikeCount,
		"comment_count": thread.CommentCount,
		"is_deleted":    thread.IsDeleted,
		"created_at":    thread.CreatedAt,
		"profile":       users,
	}

	s.logger.Info("Thread successfully created ")
	ctx.JSON(http.StatusOK, threadResponse)
}

type getThreadRequest struct {
	PublicID string `uri:"public_id"`
}

func (s *HandlersServer) GetThreadFunc(ctx *gin.Context) {
	var req getThreadRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	thread, err := s.store.GetThread(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get thread: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get thread"})
		return
	}
	s.logger.Info("Successfully get the thread")
	ctx.JSON(http.StatusOK, thread)
}

type getThreadUserRequest struct {
	PublicID string `uri:"public_id"`
}

// get thread by user
func (s *HandlersServer) GetThreadByUserFunc(ctx *gin.Context) {
	var req getThreadUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	thread, err := s.store.GetThreadUser(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get thread by user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get thread"})
		return
	}
	ctx.JSON(http.StatusOK, thread)
}

func (s *HandlersServer) GetAllThreadDetailFunc(ctx *gin.Context) {
	threads, err := s.store.GetAllThreads(ctx)
	if err != nil {
		s.logger.Error("Failed to find the all threads ", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get threads"})
		return
	}
	s.logger.Debug("Received threads from database")
	ctx.JSON(http.StatusOK, threads)
}

func (s *HandlersServer) GetAllThreadsFunc(ctx *gin.Context) {
	threads, err := s.store.GetAllThreads(ctx)
	if err != nil {
		s.logger.Error("Failed to find the all threads ", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get threads"})
		return
	}
	ctx.JSON(http.StatusOK, threads)
}

func (s *HandlersServer) GetAllThreadsByCommunitiesFunc(ctx *gin.Context) {
	var req struct {
		CommunityPublicID string `uri:"community_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	communityPublicID, err := uuid.Parse(req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	threads, err := s.store.GetAllThreadsByCommunities(ctx, communityPublicID)
	if err != nil {
		s.logger.Error("Failed to get thread by communities: ", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get threads"})
		return
	}
	s.logger.Info("Successfully get the thread")
	ctx.JSON(http.StatusOK, threads)
}

type updateThreadLikeRequest struct {
	PublicID string `uri:"public_id"`
}

func (s *HandlersServer) UpdateThreadLikeFunc(ctx *gin.Context) {
	var req updateThreadLikeRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	thread, err := s.store.UpdateThreadLike(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to update like: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like"})
		return
	}

	s.logger.Debug("Successfully update the thread ", thread)
	ctx.JSON(http.StatusOK, thread)
}

func (s *HandlersServer) UpdateThreadCommentCountFunc(ctx *gin.Context) {
	var req struct {
		PublicID string `uri:"public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	thread, err := s.store.UpdateThreadCommentCount(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to update comment count: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment count"})
		return
	}

	s.logger.Debug("Successfully updated the thread ", thread)
	ctx.JSON(http.StatusOK, thread)
}

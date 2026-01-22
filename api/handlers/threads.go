package handlers

import (
	"khelogames/core/token"
	db "khelogames/database"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createThreadRequest struct {
	CommunityPublicID string `json:"community_public_id" binding:"omitempty,uuid4"`
	Title             string `json:"title" binding:"omitempty,min=1,max=500"`
	Content           string `json:"content" binding:"omitempty,min=1,max=10000"`
	MediaType         string `json:"media_type" binding:"omitempty,oneof=image video gif link"`
	MediaURL          string `json:"media_url" binding:"omitempty,url"`
}

func (s *HandlersServer) CreateThreadFunc(ctx *gin.Context) {
	var req createThreadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	var emptyString string

	if req.Title == emptyString && req.Content == emptyString && req.MediaURL == emptyString {
		fieldErrors := map[string]string{
			"global": "Please provide any of the input",
		}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	var communityPublicID *uuid.UUID
	if req.CommunityPublicID != "" && req.CommunityPublicID != "null" {
		parsed, err := uuid.Parse(req.CommunityPublicID)
		if err != nil {
			s.logger.Error("Failed to parse community public id: ", err)
			fieldErrors := map[string]string{"community_public_id": "Invalid UUID format"}
			errorhandler.ValidationErrorResponse(ctx, fieldErrors)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create thread",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	users, err := s.store.GetProfileByUserID(ctx, thread.UserID)
	if err != nil {
		s.logger.Error("Failed to get user profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
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
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    threadResponse,
	})
}

type getThreadRequest struct {
	PublicID string `uri:"public_id" binding:"required"`
}

func (s *HandlersServer) GetThreadFunc(ctx *gin.Context) {
	var req getThreadRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	thread, err := s.store.GetThread(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get thread: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get thread",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("Successfully get the thread")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    thread,
	})
}

type getThreadUserRequest struct {
	PublicID string `uri:"public_id" binding:"required"`
}

// get thread by user
func (s *HandlersServer) GetThreadByUserFunc(ctx *gin.Context) {
	var req getThreadUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	thread, err := s.store.GetThreadUser(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get thread by user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get thread by user",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    thread,
	})
}

func (s *HandlersServer) GetAllThreadDetailFunc(ctx *gin.Context) {
	threads, err := s.store.GetAllThreads(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch threads", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Unable to fetch threads",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Even if threads is empty → OK
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    threads,
	})
}

func (s *HandlersServer) GetAllThreadsFunc(ctx *gin.Context) {
	threads, err := s.store.GetAllThreads(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch threads", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Unable to fetch threads",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Even if threads is empty → OK
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    threads,
	})
}

func (s *HandlersServer) GetAllThreadsByCommunitiesFunc(ctx *gin.Context) {
	var req struct {
		CommunityPublicID string `uri:"community_public_id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	communityPublicID, err := uuid.Parse(req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"community_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	threads, err := s.store.GetAllThreadsByCommunities(ctx, communityPublicID)
	if err != nil {
		s.logger.Error("Failed to get thread by communities: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get threads by communities",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("Successfully get the thread")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    threads,
	})
}

type updateThreadLikeRequest struct {
	PublicID string `uri:"public_id" binding:"required"`
}

func (s *HandlersServer) UpdateThreadLikeFunc(ctx *gin.Context) {
	var req updateThreadLikeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	thread, err := s.store.UpdateThreadLike(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to update like: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update like count",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Debug("Successfully update the thread ", thread)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    thread,
	})
}

func (s *HandlersServer) UpdateThreadCommentCountFunc(ctx *gin.Context) {
	var req struct {
		PublicID string `uri:"public_id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	thread, err := s.store.UpdateThreadCommentCount(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to update comment count: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update comment count",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Debug("Successfully updated the thread ", thread)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    thread,
	})
}

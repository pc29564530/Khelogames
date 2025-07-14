package handlers

import (
	db "khelogames/database"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createThreadRequest struct {
	CommunityID int32  `json:"community_id,omitempty"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MediaType   string `json:"mediaType,omitempty"`
	MediaURL    string `json:"mediaURL,omitempty"`
}

func (s *HandlersServer) CreateThreadFunc(ctx *gin.Context) {
	var req createThreadRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation")
		return
	}

	//function for uploading a image or video
	arg := db.CreateThreadParams{
		CommunityID: req.CommunityID,
		Title:       req.Title,
		Content:     req.Content,
		MediaType:   req.MediaType,
		MediaUrl:    req.MediaURL,
	}

	s.logger.Debug("received arg of create thread params: %s", arg)

	thread, err := s.store.CreateThread(ctx, arg)
	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create new thread ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit the transcation")
		return
	}

	s.logger.Info("Thread successfully created ")
	ctx.JSON(http.StatusOK, thread)
	return
}

type getThreadRequest struct {
	PublicID uuid.UUID `uri:"public_id"`
}

func (s *HandlersServer) GetThreadFunc(ctx *gin.Context) {
	var req getThreadRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	thread, err := s.store.GetThread(ctx, req.PublicID)
	if err != nil {
		s.logger.Error("Failed to get thread: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Info("Successfully get the thread")
	ctx.JSON(http.StatusOK, thread)
	return
}

type getThreadUserRequest struct {
	PublicID uuid.UUID `uri:"public_id"`
}

// get thread by user
func (s *HandlersServer) GetThreadByUserFunc(ctx *gin.Context) {
	var req getThreadUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	thread, err := s.store.GetThreadUser(ctx, req.PublicID)
	if err != nil {
		s.logger.Error("Failed to get thread by user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	ctx.JSON(http.StatusOK, thread)
	return
}

func (s *HandlersServer) GetAllThreadDetailFunc(ctx *gin.Context) {
	threads, err := s.store.GetAllThreads(ctx)
	if err != nil {
		s.logger.Error("Failed to find the all threads ", err)
		ctx.JSON(http.StatusNotFound, (err))
	}
	s.logger.Debug("Received threads from database")
	ctx.JSON(http.StatusOK, threads)
	return

}

func (s *HandlersServer) GetAllThreadsFunc(ctx *gin.Context) {
	threads, err := s.store.GetAllThreads(ctx)
	if err != nil {
		s.logger.Error("Failed to find the all threads ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	ctx.JSON(http.StatusOK, threads)
	return
}

func (s *HandlersServer) GetAllThreadsByCommunitiesFunc(ctx *gin.Context) {
	var req struct {
		CommunityPublicID uuid.UUID `uri:"community_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusInternalServerError, (err))
	}

	threads, err := s.store.GetAllThreadsByCommunities(ctx, req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Failed to get thread by communities: ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	s.logger.Info("Successfully get the thread")
	ctx.JSON(http.StatusOK, threads)
	return
}

type updateThreadLikeRequest struct {
	PublicID uuid.UUID `uri:"public_id"`
}

func (s *HandlersServer) UpdateThreadLikeFunc(ctx *gin.Context) {
	var req updateThreadLikeRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	thread, err := s.store.UpdateThreadLike(ctx, req.PublicID)
	if err != nil {
		s.logger.Error("Failed to update like: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug("Successfully update the thread ", thread)
	ctx.JSON(http.StatusOK, thread)
	return
}

func (s *HandlersServer) UpdateThreadCommentCountFunc(ctx *gin.Context) {
	var req struct {
		PublicID uuid.UUID `uri:"public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	thread, err := s.store.UpdateThreadCommentCount(ctx, req.PublicID)
	if err != nil {
		s.logger.Error("Failed to update like: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug("Successfully update the thread ", thread)
	ctx.JSON(http.StatusOK, thread)
	return
}

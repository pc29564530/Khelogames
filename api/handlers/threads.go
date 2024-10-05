package handlers

import (
	"encoding/base64"
	db "khelogames/database"

	"khelogames/pkg"
	"khelogames/token"
	util "khelogames/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type createThreadRequest struct {
	CommunitiesName string `json:"communities_name,omitempty"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	MediaType       string `json:"mediaType,omitempty"`
	MediaURL        string `json:"mediaURL,omitempty"`
	LikeCount       int64  `json:"likeCount"`
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

	saveImageStruct := util.NewSaveImageStruct(s.logger)

	var path string
	if req.MediaType != "" {
		b64data := req.MediaURL[strings.IndexByte(req.MediaURL, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			s.logger.Error("Failed to decode string", err)
			return
		}

		path, err = saveImageStruct.SaveImageToFile(data, req.MediaType)
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to save image to file ", err)
			return
		}
	}

	//function for uploading a image or video
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateThreadParams{
		Username:        authPayload.Username,
		CommunitiesName: req.CommunitiesName,
		Title:           req.Title,
		Content:         req.Content,
		MediaType:       req.MediaType,
		MediaUrl:        path,
		LikeCount:       0,
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
	ID int64 `uri:"id"`
}

func (s *HandlersServer) GetThreadFunc(ctx *gin.Context) {
	var req getThreadRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	thread, err := s.store.GetThread(ctx, req.ID)
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
	Username string `uri:"username"`
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

	thread, err := s.store.GetThreadUser(ctx, req.Username)
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
		return
	}
	s.logger.Debug("Received threads from database")
	var threadsDetails []map[string]interface{}

	for _, thread := range threads {
		profile, err := s.store.GetProfile(ctx, thread.Username)
		if err != nil {
			s.logger.Error("Failed to find the profile ", err)
			return
		}
		var displayText string
		if profile.AvatarUrl == "" {
			displayText = strings.ToUpper(string(profile.FullName[0]))
		}

		threadsDetail := map[string]interface{}{
			"id":               thread.ID,
			"username":         thread.Username,
			"communities_name": thread.CommunitiesName,
			"title":            thread.Title,
			"content":          thread.Content,
			"media_type":       thread.MediaType,
			"media_url":        thread.MediaUrl,
			"like_count":       thread.LikeCount,
			"full_name":        profile.FullName,
			"avatar_url":       profile.AvatarUrl,
			"display_text":     displayText,
			"created_at":       thread.CreatedAt,
		}
		threadsDetails = append(threadsDetails, threadsDetail)
	}
	ctx.JSON(http.StatusOK, threadsDetails)
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

type getThreadsByCommunitiesRequest struct {
	CommunitiesName string `uri:"communities_name"`
}

func (s *HandlersServer) GetAllThreadsByCommunityDetailsFunc(ctx *gin.Context) {
	var req getThreadsByCommunitiesRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusInternalServerError, (err))
	}

	s.logger.Info("get community name: %s", req.CommunitiesName)

	threads, err := s.store.GetAllThreadsByCommunities(ctx, req.CommunitiesName)
	if err != nil {
		s.logger.Error("Failed to get thread by communities: ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}

	s.logger.Info("Received threads from database", threads)
	var threadsDetails []map[string]interface{}

	for _, thread := range threads {
		profile, err := s.store.GetProfile(ctx, thread.Username)
		if err != nil {
			s.logger.Error("Failed to find the profile ", err)
			return
		}

		var displayText string
		if profile.AvatarUrl == "" {
			displayText = strings.ToUpper(string(profile.FullName[0]))
		}

		threadsDetail := map[string]interface{}{
			"id":               thread.ID,
			"username":         thread.Username,
			"communities_name": thread.CommunitiesName,
			"title":            thread.Title,
			"content":          thread.Content,
			"media_type":       thread.MediaType,
			"media_url":        thread.MediaUrl,
			"like_count":       thread.LikeCount,
			"full_name":        profile.FullName,
			"avatar_url":       profile.AvatarUrl,
			"display_text":     displayText,
			"created_at":       thread.CreatedAt,
		}
		threadsDetails = append(threadsDetails, threadsDetail)
	}
	ctx.JSON(http.StatusOK, threadsDetails)
	return
}

func (s *HandlersServer) GetAllThreadsByCommunitiesFunc(ctx *gin.Context) {
	var req getThreadsByCommunitiesRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusInternalServerError, (err))
	}

	threads, err := s.store.GetAllThreadsByCommunities(ctx, req.CommunitiesName)
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
	LikeCount int64 `json:"like_count"`
	ID        int64 `json:"id"`
}

func (s *HandlersServer) UpdateThreadLikeFunc(ctx *gin.Context) {
	var req updateThreadLikeRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	arg := db.UpdateThreadLikeParams{
		LikeCount: req.LikeCount,
		ID:        req.ID,
	}

	thread, err := s.store.UpdateThreadLike(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update like: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug("Successfully update the thread ", thread)
	ctx.JSON(http.StatusOK, thread)
	return
}

type threadByThreadIdRequest struct {
	ID int64 `uri:"id"`
}

func (s *HandlersServer) GetThreadByThreadIDFunc(ctx *gin.Context) {
	var req threadByThreadIdRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	thread, err := s.store.GetThreadByThreadID(ctx, req.ID)
	if err != nil {
		s.logger.Error("Failed to get thread by thread id ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	s.logger.Info("Successfully get thread by thread id ")
	ctx.JSON(http.StatusAccepted, thread)
	return
}

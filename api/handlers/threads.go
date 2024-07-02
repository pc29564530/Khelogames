package handlers

import (
	"encoding/base64"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	util "khelogames/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ThreadServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewThreadServer(store *db.Store, logger *logger.Logger) *ThreadServer {
	return &ThreadServer{store: store, logger: logger}
}

type createThreadRequest struct {
	CommunitiesName string `json:"communities_name,omitempty"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	MediaType       string `json:"mediaType,omitempty"`
	MediaURL        string `json:"mediaURL,omitempty"`
	LikeCount       int64  `json:"likeCount"`
}

func (s *ThreadServer) CreateThreadFunc(ctx *gin.Context) {
	var req createThreadRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	var path string
	if req.MediaType != "" {
		b64data := req.MediaURL[strings.IndexByte(req.MediaURL, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			fmt.Errorf("Failed to decode string", err)
			return
		}

		path, err = util.SaveImageToFile(data, req.MediaType)
		if err != nil {
			fmt.Errorf("Failed to save image to file ", err)
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

	thread, err := s.store.CreateThread(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create new thread ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	fmt.Println("Thread successfully created ")
	ctx.JSON(http.StatusOK, thread)
	return
}

type getThreadRequest struct {
	ID int64 `uri:"id"`
}

func (s *ThreadServer) GetThreadFunc(ctx *gin.Context) {
	var req getThreadRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	thread, err := s.store.GetThread(ctx, req.ID)
	if err != nil {
		fmt.Errorf("Failed to get thread: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	fmt.Println("Successfully get the thread")
	ctx.JSON(http.StatusOK, thread)
	return
}

type getThreadUserRequest struct {
	Username string `uri:"username"`
}

// get thread by user
func (s *ThreadServer) GetThreadByUserFunc(ctx *gin.Context) {
	var req getThreadUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	thread, err := s.store.GetThreadUser(ctx, req.Username)
	if err != nil {
		fmt.Errorf("Failed to get thread by user: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	ctx.JSON(http.StatusOK, thread)
	return
}

func (s *ThreadServer) GetAllThreadsFunc(ctx *gin.Context) {
	threads, err := s.store.GetAllThreads(ctx)
	if err != nil {
		fmt.Errorf("Failed to find the all threads ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	ctx.JSON(http.StatusOK, threads)
	return
}

type getThreadsByCommunitiesRequest struct {
	CommunitiesName string `uri:"communities_name"`
}

func (s *ThreadServer) GetAllThreadsByCommunitiesFunc(ctx *gin.Context) {
	var req getThreadsByCommunitiesRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind", err)
		ctx.JSON(http.StatusInternalServerError, (err))
	}

	threads, err := s.store.GetAllThreadsByCommunities(ctx, req.CommunitiesName)
	if err != nil {
		fmt.Errorf("Failed to get thread by communities: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	fmt.Println("Successfully get the thread")
	ctx.JSON(http.StatusOK, threads)
	return
}

type updateThreadLikeRequest struct {
	LikeCount int64 `json:"like_count"`
	ID        int64 `json:"id"`
}

func (s *ThreadServer) UpdateThreadLikeFunc(ctx *gin.Context) {
	var req updateThreadLikeRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	arg := db.UpdateThreadLikeParams{
		LikeCount: req.LikeCount,
		ID:        req.ID,
	}

	thread, err := s.store.UpdateThreadLike(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update like: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	fmt.Println("Successfully update the thread ", err)
	ctx.JSON(http.StatusOK, thread)
	return
}

type threadByThreadIdRequest struct {
	ID int64 `uri:"id"`
}

func (s *ThreadServer) GetThreadByThreadIDFunc(ctx *gin.Context) {
	var req threadByThreadIdRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	thread, err := s.store.GetThreadByThreadID(ctx, req.ID)
	if err != nil {
		fmt.Errorf("Failed to get thread by thread id ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	fmt.Println("Successfully get thread by thread id ")
	ctx.JSON(http.StatusAccepted, thread)
	return
}

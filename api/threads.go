package api

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func saveImageToFile(data []byte, mediaType string) (string, error) {
	randomString, err := generateRandomString(12)
	if err != nil {
		fmt.Printf("Error generating random string: %v\n", err)
		return "", err
	}

	var mediaFolder string
	switch mediaType {
	case "image":
		mediaFolder = "images"
	case "video":
		mediaFolder = "videos"
	default:
		return "", fmt.Errorf("unsupported media type for inserting in mediaFolder: %s", mediaType)
	}

	filePath := fmt.Sprintf("/Users/pawan/database/Khelogames/%s/%s", mediaFolder, randomString)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	path := convertLocalPathToURL(filePath, mediaFolder)
	return path, nil
}

func generateRandomString(length int) (string, error) {
	if length%2 != 0 {
		return "", fmt.Errorf("length must be even for generating hex string")
	}

	randomBytes := make([]byte, length/2)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}

func convertLocalPathToURL(localPath string, mediaFolder string) string {
	baseURL := fmt.Sprintf("http://10.0.2.2:8080/%s/", mediaFolder)
	imagePath := baseURL + strings.TrimPrefix(localPath, fmt.Sprintf("/Users/pawan/database/Khelogames/%s/", mediaFolder))
	filePath := imagePath
	return filePath
}

type createThreadRequest struct {
	CommunitiesName string `json:"communities_name,omitempty"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	MediaType       string `json:"mediaType,omitempty"`
	MediaURL        string `json:"mediaURL,omitempty"`
	LikeCount       int64  `json:"likeCount"`
}

func (server *Server) createThread(ctx *gin.Context) {
	var req createThreadRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println("Hello line no 99")

	var path string
	if req.MediaType != "" {
		b64data := req.MediaURL[strings.IndexByte(req.MediaURL, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			fmt.Println("unable to decode :", err)
			return
		}

		fmt.Println("hello line no 109")

		path, err = saveImageToFile(data, req.MediaType)
		if err != nil {
			fmt.Println("unable to create a file")
			return
		}
	}

	//function for uploading a image or video
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateThreadParams{
		Username:        authPayload.Username,
		CommunitiesName: req.CommunitiesName,
		Title:           req.Title,
		Content:         req.Content,
		MediaType:       req.MediaType,
		MediaUrl:        path,
		LikeCount:       0,
	}

	thread, err := server.store.CreateThread(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, thread)
	return
}

type getThreadRequest struct {
	ID int64 `uri:"id"`
}

func (server *Server) getThread(ctx *gin.Context) {
	var req getThreadRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	thread, err := server.store.GetThread(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, thread)
	return
}

func (server *Server) getAllThreads(ctx *gin.Context) {
	threads, err := server.store.GetAllThreads(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, threads)
	return
}

type getThreadsByCommunitiesRequest struct {
	CommunitiesName string `uri:"communities_name"`
}

func (server *Server) getAllThreadsByCommunities(ctx *gin.Context) {
	var req getThreadsByCommunitiesRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	threads, err := server.store.GetAllThreadsByCommunities(ctx, req.CommunitiesName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	fmt.Println("Thread Data: ", threads)
	ctx.JSON(http.StatusOK, threads)
	return
}

type updateThreadLikeRequest struct {
	LikeCount int64 `json:"like_count"`
	ID        int64 `json:"id"`
}

func (server *Server) updateThreadLike(ctx *gin.Context) {
	var req updateThreadLikeRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateThreadLikeParams{
		LikeCount: req.LikeCount,
		ID:        req.ID,
	}

	thread, err := server.store.UpdateThreadLike(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, thread)
	return
}

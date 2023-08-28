package api

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
	"os"
	"strings"
)

type createThreadRequest struct {
	CommunitiesName string `json:"communities_name,omitempty"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	MediaType       string `json:"mediaType,omitempty"`
	MediaURL        string `json:"mediaURL,omitempty"`
	LikeCount       int64  `json:"likeCount"`
}

func (server *Server) createThread(ctx *gin.Context) {
	fmt.Println("line no 66")
	var req createThreadRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	b64data := req.MediaURL[strings.IndexByte(req.MediaURL, ',')+1:]

	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Println("uanble to decode :", err)
		return
	}

	path, err := saveImageToFile(data)
	if err != nil {
		fmt.Println("uanble to create a file")
		return
	}

	fmt.Println(path)

	//function for uplo	ading a image or video
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
	CommunitiesName string `json:"communities_name"`
}

func (server *Server) getAllThreadsByCommunities(ctx *gin.Context) {
	var req getThreadsByCommunitiesRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	threads, err := server.store.GetAllThreadsByCommunities(ctx, req.CommunitiesName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, threads)
	return
}

type updateThreadLikeRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) updateThreadLike(ctx *gin.Context) {
	var req updateThreadLikeRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	thread, err := server.store.UpdateThreadLike(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, thread)
	return

}

func (server *Server) Uploads(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	fmt.Println(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println("Hello IUndia what are you doing ")

}

//func decodeBase64Image(base64Str string) ([]byte, error) {
//	fmt.Println("line no 102")
//	data, err := base64.StdEncoding.DecodeString(base64Str)
//	fmt.Println(string(data))
//	fmt.Println(err)
//	if err != nil {
//		return nil, err
//	}
//
//	return data, nil
//}

func saveImageToFile(data []byte) (string, error) {
	filePath := "/home/pawan/projects/golang-project/Khelogames/images/image.jpg"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

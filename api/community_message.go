package api

import (
	"encoding/base64"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type createCommunityMessageRequest struct {
	CommunityName  string `json:"community_name"`
	SenderUsername string `json:"sender_username"`
	Content        string `json:"content"`
}

func (server *Server) createCommunityMessage(ctx *gin.Context) {
	var req createCommunityMessageRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateCommunityMessageParams{
		CommunityName:  req.CommunityName,
		SenderUsername: authPayload.Username,
		Content:        req.Content,
	}

	response, err := server.store.CreateCommunityMessage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) createUploadMedia(ctx *gin.Context) {

	r := ctx.Request
	if err := r.ParseMultipartForm(40 << 30); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	mediaUrl := ctx.Request.FormValue("media_url")
	mediaType := ctx.Request.FormValue("media_type")

	var path string
	if mediaType != "" {
		b64data := mediaUrl[strings.IndexByte(mediaUrl, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			fmt.Println("unable to decode :", err)
			return
		}

		path, err = saveImageToFile(data, mediaType)
		if err != nil {
			fmt.Println("unable to create a file")
			return
		}
	}

	arg := db.CreateUploadMediaParams{
		MediaUrl:  path,
		MediaType: mediaType,
	}

	response, err := server.store.CreateUploadMedia(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type createMessageMediaRequest struct {
	MessageID int64 `json:"message_id"`
	MediaID   int64 `json:"media_id"`
}

func (server *Server) createMessageMedia(ctx *gin.Context) {
	var req createMessageMediaRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, errorResponse(err))
		return
	}

	arg := db.CreateMessageMediaParams{
		MessageID: req.MessageID,
		MediaID:   req.MediaID,
	}

	response, err := server.store.CreateMessageMedia(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getCommuntiyMessage(ctx *gin.Context) {
	response, err := server.store.GetCommuntiyMessage(ctx)
	fmt.Println("Error: ", err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getCommunityByMessage(ctx *gin.Context) {
	response, err := server.store.GetCommunityByMessage(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println(response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

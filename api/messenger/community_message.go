package messenger

import (
	"encoding/base64"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CommunityMessageServer struct {
	store     *db.Store
	logger    *logger.Logger
	broadcast chan []byte
}

func NewCommunityMessageSever(store *db.Store, logger *logger.Logger, broadcast chan []byte) *CommunityMessageServer {
	return &CommunityMessageServer{store: store, logger: logger, broadcast: broadcast}
}

type createCommunityMessageRequest struct {
	CommunityName  string `json:"community_name"`
	SenderUsername string `json:"sender_username"`
	Content        string `json:"content"`
}

func (s *CommunityMessageServer) CreateCommunityMessageFunc(ctx *gin.Context) {
	var req createCommunityMessageRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateCommunityMessageParams{
		CommunityName:  req.CommunityName,
		SenderUsername: authPayload.Username,
		Content:        req.Content,
	}

	response, err := s.store.CreateCommunityMessage(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create community message: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *CommunityMessageServer) CreateUploadMediaFunc(ctx *gin.Context) {

	r := ctx.Request
	if err := r.ParseMultipartForm(40 << 30); err != nil {
		fmt.Errorf("Failed to parse multipart form create upload media: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	mediaUrl := ctx.Request.FormValue("media_url")
	mediaType := ctx.Request.FormValue("media_type")

	var path string
	if mediaType != "" {
		b64data := mediaUrl[strings.IndexByte(mediaUrl, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			fmt.Errorf("Failed to decode string: %v", err)
			return
		}

		path, err = util.SaveImageToFile(data, mediaType)
		if err != nil {
			fmt.Errorf("Failed to save image to file: %v", err)
			return
		}
	}

	arg := db.CreateUploadMediaParams{
		MediaUrl:  path,
		MediaType: mediaType,
	}

	response, err := s.store.CreateUploadMedia(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create upload media: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type createMessageMediaRequest struct {
	MessageID int64 `json:"message_id"`
	MediaID   int64 `json:"media_id"`
}

func (s *CommunityMessageServer) CreateMessageMediaFunc(ctx *gin.Context) {
	var req createMessageMediaRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusBadGateway, (err))
		return
	}

	arg := db.CreateMessageMediaParams{
		MessageID: req.MessageID,
		MediaID:   req.MediaID,
	}

	response, err := s.store.CreateMessageMedia(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create message media: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *CommunityMessageServer) GetCommuntiyMessageFunc(ctx *gin.Context) {
	response, err := s.store.GetCommuntiyMessage(ctx)
	if err != nil {
		fmt.Errorf("Failed to get community message: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *CommunityMessageServer) GetCommunityByMessageFunc(ctx *gin.Context) {
	response, err := s.store.GetCommunityByMessage(ctx)
	if err != nil {
		fmt.Errorf("Failed to get community by message: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

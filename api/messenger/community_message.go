package messenger

import (
	"encoding/base64"
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

func NewCommunityMessageServer(store *db.Store, logger *logger.Logger, broadcast chan []byte) *CommunityMessageServer {
	return &CommunityMessageServer{store: store, logger: logger, broadcast: broadcast}
}

type createCommunityMessageRequest struct {
	CommunityName  string `json:"community_name"`
	SenderUsername string `json:"sender_username"`
	Content        string `json:"content"`
}

func (s *CommunityMessageServer) CreateCommunityMessageFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create community message")

	var req createCommunityMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error("Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: %v", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateCommunityMessageParams{
		CommunityName:  req.CommunityName,
		SenderUsername: authPayload.Username,
		Content:        req.Content,
	}

	s.logger.Debug("Create community message params: %v", arg)

	response, err := s.store.CreateCommunityMessage(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create community message: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Info("Successfully created community message")
	ctx.JSON(http.StatusAccepted, response)
}

func (s *CommunityMessageServer) CreateUploadMediaFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create upload media")

	r := ctx.Request
	if err := r.ParseMultipartForm(40 << 20); err != nil {
		s.logger.Error("Failed to parse multipart form: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	mediaUrl := ctx.Request.FormValue("media_url")
	mediaType := ctx.Request.FormValue("media_type")

	s.logger.Debug("Received create upload media params")
	saveImageStruct := util.NewSaveImageStruct(s.logger)

	var path string
	if mediaType != "" {
		b64data := mediaUrl[strings.IndexByte(mediaUrl, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			s.logger.Error("Failed to decode base64 string: %v", err)
			return
		}

		path, err = saveImageStruct.SaveImageToFile(data, mediaType)
		if err != nil {
			s.logger.Error("Failed to save image to file: %v", err)
			return
		}
		s.logger.Debug("Image saved successfully at %s", path)
	}

	arg := db.CreateUploadMediaParams{
		MediaUrl:  path,
		MediaType: mediaType,
	}

	s.logger.Debug("Create upload media params: %v", arg)

	response, err := s.store.CreateUploadMedia(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create upload media: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Info("Successfully created upload media")
	ctx.JSON(http.StatusAccepted, response)
}

type createMessageMediaRequest struct {
	MessageID int64 `json:"message_id"`
	MediaID   int64 `json:"media_id"`
}

func (s *CommunityMessageServer) CreateMessageMediaFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create message media")

	var req createMessageMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error("Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: %v", req)

	arg := db.CreateMessageMediaParams{
		MessageID: req.MessageID,
		MediaID:   req.MediaID,
	}

	s.logger.Debug("Create message media params: %v", arg)

	response, err := s.store.CreateMessageMedia(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create message media: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully created message media")
	ctx.JSON(http.StatusAccepted, response)
}

func (s *CommunityMessageServer) GetCommunityMessageFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get community message")

	response, err := s.store.GetCommuntiyMessage(ctx) //spelling mistake
	if err != nil {
		s.logger.Error("Failed to get community message: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved community message")
	ctx.JSON(http.StatusOK, response)
}

func (s *CommunityMessageServer) GetCommunityByMessageFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get community by message")

	response, err := s.store.GetCommunityByMessage(ctx)
	if err != nil {
		s.logger.Error("Failed to get community by message: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved community by message")
	ctx.JSON(http.StatusOK, response)
}

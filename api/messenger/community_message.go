package messenger

import (
	"encoding/base64"
	db "khelogames/db/sqlc"
	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type createCommunityMessageRequest struct {
	CommunityName  string `json:"community_name"`
	SenderUsername string `json:"sender_username"`
	Content        string `json:"content"`
}

func (s *MessageServer) CreateCommunityMessageFunc(ctx *gin.Context) {

	var req createCommunityMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error("Failed to bind JSON: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateCommunityMessageParams{
		CommunityName:  req.CommunityName,
		SenderUsername: authPayload.Username,
		Content:        req.Content,
	}

	s.logger.Debug("Create community message params: ", arg)

	response, err := s.store.CreateCommunityMessage(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create community message: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Info("Successfully created community message")
	ctx.JSON(http.StatusAccepted, response)
}

func (s *MessageServer) CreateUploadMediaFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create upload media")

	r := ctx.Request
	if err := r.ParseMultipartForm(40 << 20); err != nil {
		s.logger.Error("Failed to parse multipart form: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin the transcation")
	}

	defer tx.Rollback()

	mediaUrl := ctx.Request.FormValue("media_url")
	mediaType := ctx.Request.FormValue("media_type")

	s.logger.Debug("Received create upload media params")
	saveImageStruct := util.NewSaveImageStruct(s.logger)

	var path string
	if mediaType != "" {
		b64data := mediaUrl[strings.IndexByte(mediaUrl, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			s.logger.Error("Failed to decode base64 string: ", err)
			return
		}

		path, err = saveImageStruct.SaveImageToFile(data, mediaType)
		if err != nil {
			s.logger.Error("Failed to save image to file: ", err)
			return
		}
		s.logger.Debug("Image saved successfully at %s", path)
	}

	arg := db.CreateUploadMediaParams{
		MediaUrl:  path,
		MediaType: mediaType,
	}

	s.logger.Debug("Create upload media params: ", arg)

	response, err := s.store.CreateUploadMedia(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create upload media: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	s.logger.Info("Successfully created upload media")
	ctx.JSON(http.StatusAccepted, response)
}

type createMessageMediaRequest struct {
	MessageID int64 `json:"message_id"`
	MediaID   int64 `json:"media_id"`
}

func (s *MessageServer) CreateMessageMediaFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create message media")

	var req createMessageMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error("Failed to bind JSON: ", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	s.logger.Debug("Request JSON bind successful: ", req)

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation: ", err)
		return
	}

	defer tx.Rollback()

	arg := db.CreateMessageMediaParams{
		MessageID: req.MessageID,
		MediaID:   req.MediaID,
	}

	s.logger.Debug("Create message media params: ", arg)

	response, err := s.store.CreateMessageMedia(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create message media: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	s.logger.Info("Successfully created message media")
	ctx.JSON(http.StatusAccepted, response)
}

func (s *MessageServer) GetCommunityMessageFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get community message")

	response, err := s.store.GetCommuntiyMessage(ctx) //spelling mistake
	if err != nil {
		s.logger.Error("Failed to get community message: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved community message")
	ctx.JSON(http.StatusOK, response)
}

func (s *MessageServer) GetCommunityByMessageFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get community by message")

	response, err := s.store.GetCommunityByMessage(ctx)
	if err != nil {
		s.logger.Error("Failed to get community by message: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved community by message")
	ctx.JSON(http.StatusOK, response)
}

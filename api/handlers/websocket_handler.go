package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ampq "github.com/rabbitmq/amqp091-go"
)

type WebSocketHandlerImpl struct {
	store      *db.Store
	tokenMaker token.Maker
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	rabbitChan *ampq.Channel
	mutex      sync.Mutex
	logger     *logger.Logger
}

func NewWebSocketHandler(store *db.Store, tokenMaker token.Maker, rabbitChan *ampq.Channel, logger *logger.Logger) *WebSocketHandlerImpl {
	return &WebSocketHandlerImpl{
		store:      store,
		tokenMaker: tokenMaker,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		rabbitChan: rabbitChan,
		logger:     logger,
	}
}

func (h *WebSocketHandlerImpl) HandleWebSocket(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	auth := strings.Split(authHeader, " ")

	if len(auth) == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	_, err := h.tokenMaker.VerifyToken(auth[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	h.clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			delete(h.clients, conn)
			break
		}

		var message map[string]string
		err = json.Unmarshal(msg, &message)
		if err != nil {
			fmt.Print("unable to unmarshal msg ", err)
			return
		}

		err = h.rabbitChan.Publish(
			"",
			"message",
			false,
			false,
			ampq.Publishing{
				ContentType: "application/json",
				Body:        msg,
			},
		)

		authToken := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
		b64data := message["media_url"][strings.IndexByte(message["media_url"], ',')+1:]
		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			fmt.Print("unable to decode :", err)
			return
		}
		mediaType := "image"
		path, err := util.SaveImageToFile(data, mediaType)
		if err != nil {
			fmt.Print("unable to create a file")
			return
		}

		arg := db.CreateNewMessageParams{
			Content:          message["content"],
			IsSeen:           false,
			SenderUsername:   authToken.Username,
			ReceiverUsername: message["receiver_username"],
			MediaUrl:         path,
			MediaType:        message["media_type"],
		}

		_, err = h.store.CreateNewMessage(ctx, arg)
		if err != nil {
			fmt.Print("unable to store new message: ", err)
			return
		}

		h.broadcast <- msg
	}
}

package messenger

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

func StartRabbitMQ(config util.Config) (*ampq.Connection, *ampq.Channel, error) {
	rabbitConn, err := ampq.Dial(config.RabbitSource)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to RabbitMQ :%w", err)
	}
	rabbitChan, err := rabbitConn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open RabbitMQ channel :%w", err)
	}
	return rabbitConn, rabbitChan, nil
}

func (s *WebSocketHandlerImpl) StartWebSocketHub() {
	for {
		select {
		case message := <-s.broadcast:
			s.mutex.Lock()
			for client := range s.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					delete(s.clients, client)
					client.Close()
				}
			}
			s.mutex.Unlock()
		}
	}
}

func NewWebSocketHandler(store *db.Store, tokenMaker token.Maker, clients map[*websocket.Conn]bool, broadcast chan []byte, upgrader websocket.Upgrader, rabbitChan *ampq.Channel, logger *logger.Logger) *WebSocketHandlerImpl {
	return &WebSocketHandlerImpl{
		store:      store,
		tokenMaker: tokenMaker,
		upgrader:   upgrader,
		clients:    clients,
		broadcast:  broadcast,
		rabbitChan: rabbitChan,
		logger:     logger,
	}
}

func (h *WebSocketHandlerImpl) HandleWebSocket(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	auth := strings.Split(authHeader, " ")

	if len(auth) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	payload, err := h.tokenMaker.VerifyToken(auth[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		h.logger.Errorf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()

	defer func() {
		h.mutex.Lock()
		delete(h.clients, conn)
		h.mutex.Unlock()
	}()

	h.logger.Infof("WebSocket connection established for user %s", payload.Username)

	for {
		fmt.Println("Line no 196")
		_, msg, err := conn.ReadMessage()
		fmt.Println("Error: ", err)
		fmt.Println("Line no 108: ", string(msg))
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

		fmt.Println("WebSocket: ", message)
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

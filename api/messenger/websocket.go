package messenger

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	db "khelogames/db/sqlc"

	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ampq "github.com/rabbitmq/amqp091-go"
)

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

func (s *MessageServer) StartWebSocketHub() {
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

func (h *MessageServer) HandleWebSocket(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	auth := strings.Split(authHeader, " ")

	if len(auth) != 2 {
		h.logger.Error("no token provided")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	payload, err := h.tokenMaker.VerifyToken(auth[1])
	if err != nil {
		h.logger.Debug("unable to get valid token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	h.logger.Debug("payload of verify token: %v", payload)

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		h.logger.Errorf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	h.logger.Debug("upgrade the websocket: %v", conn)

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
		if err != nil {
			h.logger.Error("unable to read message: %v", err)
			delete(h.clients, conn)
			break
		}

		h.logger.Debug("successfully read message: %v", msg)

		var message map[string]string
		err = json.Unmarshal(msg, &message)
		if err != nil {
			h.logger.Error("unable to unmarshal msg ", err)
			return
		}

		h.logger.Debug("unmarshal message successfully ", message)
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

		if err != nil {
			h.logger.Error("unable to publish message to rabbitchannel: %v", err)
			return
		}

		authToken := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
		b64data := message["media_url"][strings.IndexByte(message["media_url"], ',')+1:]
		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			h.logger.Error("unable to decode string: %v", err)
			return
		}
		saveImageStruct := util.NewSaveImageStruct(h.logger)
		mediaType := "image"
		path, err := saveImageStruct.SaveImageToFile(data, mediaType)
		if err != nil {
			h.logger.Error("unable to create a file")
			return
		}

		h.logger.Debug("image path successfully created: %v", path)
		h.logger.Info("successfully created image path")

		arg := db.CreateNewMessageParams{
			Content:          message["content"],
			IsSeen:           false,
			SenderUsername:   authToken.Username,
			ReceiverUsername: message["receiver_username"],
			MediaUrl:         path,
			MediaType:        message["media_type"],
		}

		h.logger.Debug("create new message params: %v", arg)

		_, err = h.store.CreateNewMessage(ctx, arg)
		if err != nil {
			h.logger.Error("unable to store new message: ", err)
			return
		}

		h.logger.Info("successfully created a new message")

		h.broadcast <- msg

		h.logger.Debug("Successfully broad cast message")
	}
}

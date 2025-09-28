package messenger

import (
	"encoding/json"
	"fmt"
	db "khelogames/database"

	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		h.logger.Debug("unable to get valid token: ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	h.logger.Debug("payload of verify token: ", payload)

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade to WebSocket: ", err)
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

	h.logger.Infof("WebSocket connection established for user %s", payload.PublicID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.logger.Error("unable to read message: ", err)
			delete(h.clients, conn)
			break
		}

		h.logger.Debug("successfully read message: ", msg)

		var message map[string]interface{}
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
			h.logger.Error("unable to publish message to rabbitchannel: ", err)
			return
		}

		tx, err := h.store.BeginTx(ctx)
		if err != nil {
			h.logger.Error("Failed to begin transcation: ", err)
			return
		}

		defer tx.Rollback()

		authToken := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
		receiverPublicID, err := uuid.Parse(message["receiver_public_id"].(string))
		if err != nil {
			h.logger.Error("Failed to parse to uuid: ", err)
			return
		}

		sentAtStr := message["sent_at"].(string) // e.g. "2025-09-14T12:30:00Z"
		sentAt, err := time.Parse(time.RFC3339, sentAtStr)
		if err != nil {
			h.logger.Warn("invalid sent_at format, using now() instead")
			sentAt = time.Now()
		}
				
		arg := db.CreateNewMessageParams{
			SenderID:   authToken.PublicID,
			ReceiverID: receiverPublicID,
			Content:    message["content"].(string),
			MediaUrl:   message["media_url"].(string),
			MediaType:  message["media_type"].(string),
			SentAt: sentAt,
		}

		h.logger.Debug("create new message params: ", arg)

		_, err = h.store.CreateNewMessage(ctx, arg)
		if err != nil {
			h.logger.Error("unable to store new message: ", err)
			return
		}

		h.logger.Info("successfully created a new message")

		h.broadcast <- msg

		h.logger.Debug("Successfully broad cast message")

		err = tx.Commit()
		if err != nil {
			h.logger.Error("Failed to commit the transcation: ", err)
			return
		}
	}
}

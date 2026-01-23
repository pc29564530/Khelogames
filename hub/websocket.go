package hub

import (
	"encoding/json"
	"fmt"
	"khelogames/core/token"
	"khelogames/database"
	"khelogames/pkg"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Hub) HandleWebSocket(ctx *gin.Context) {
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

	client := h.AddClient(conn, payload.PublicID)
	if err != nil {
		h.logger.Error("Failed to add client to WebSocket: ", err)
		return
	}
	defer h.RemoveClient(client)

	// h.mu.Lock()
	// h.clients[conn] = true
	// h.mutex.Unlock()

	// defer func() {
	// 	h.mutex.Lock()
	// 	delete(h.clients, conn)
	// 	h.mutex.Unlock()
	// }()

	h.logger.Infof("WebSocket connection established for user %s", payload.PublicID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.logger.Error("unable to read message: ", err)
			// delete(h.Clients)
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
		payload, ok := message["payload"].(map[string]interface{})
		if !ok {
			h.logger.Error("invalid payload format")
			continue
		}

		switch message["type"].(string) {
		case "SUBSCRIBE":
			switch message["category"].(string) {
			case "CHAT":
				h.SubscribeClient(client, payload["profile_public_id"].(string))
			case "MATCH":
				h.SubscribeClient(client, payload["match_public_id"].(string))
			}
		case "CREATE_MESSAGE":
			h.CreateMessage(ctx, msg, message["payload"].(map[string]interface{}))
		}
	}
}

func (h *Hub) SubscribeClient(client *Client, topic string) {
	// fmt.Println("Client Line no 94: ", client)
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subscriber[topic] == nil {
		h.subscriber[topic] = make(map[*Client]bool)
	}
	h.subscriber[topic][client] = true
	h.logger.Infof("Client %s subscribed to %s", client.UserPublicID, topic)
}

func (h *Hub) CreateMessage(ctx *gin.Context, msg []byte, message map[string]interface{}) {
	// err := h.rabbitChan.PublishWithContext(
	// 	ctx,
	// 	"",
	// 	"message",
	// 	false,
	// 	false,
	// 	ampq.Publishing{
	// 		ContentType: "application/json",
	// 		Body:        msg,
	// 	},
	// )

	// if err != nil {
	// 	h.logger.Error("unable to publish message to rabbitchannel: ", err)
	// 	return
	// }

	authToken := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	fmt.Println("Receiver Id: ", message["receiver_public_id"])
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

	arg := database.CreateNewMessageParams{
		SenderID:   authToken.PublicID,
		ReceiverID: receiverPublicID,
		Content:    message["content"].(string),
		MediaUrl:   message["media_url"].(string),
		MediaType:  message["media_type"].(string),
		SentAt:     sentAt,
	}

	h.logger.Debug("create new message params: ", arg)

	msgData, err := h.store.CreateNewMessage(ctx, arg)
	if err != nil {
		h.logger.Error("unable to store new message: ", err)
		return
	}

	sender, err := h.store.GetProfile(ctx, authToken.PublicID)
	if err != nil {
		h.logger.Error("Failed to get profile by public id: ", err)
		return
	}

	receiver, err := h.store.GetProfileByPublicID(ctx, receiverPublicID)
	if err != nil {
		h.logger.Error("Failed to get profile by public id: ", err)
		return
	}

	data := map[string]interface{}{
		"id":        msgData.ID,
		"public_id": msgData.PublicID,
		"sender": map[string]interface{}{
			"public_id":  sender.PublicID,
			"user_id":    sender.UserID,
			"full_name":  sender.FullName,
			"avatar_url": sender.AvatarUrl,
			"bio":        sender.Bio,
			"username":   sender.Username,
			"created_at": sender.CreatedAT,
		},
		"receiver": map[string]interface{}{
			"public_id":  receiver.PublicID,
			"user_id":    receiver.UserID,
			"full_name":  receiver.FullName,
			"avatar_url": receiver.AvatarUrl,
			"bio":        receiver.Bio,
			"username":   receiver.Username,
			"created_at": receiver.CreatedAT,
		},
		"content":      msgData.Content,
		"media_url":    msgData.MediaUrl,
		"media_type":   msgData.MediaType,
		"is_delivered": msgData.IsDelivered,
		"is_seen":      msgData.IsSeen,
		"is_deleted":   msgData.IsDeleted,
		"created_at":   msgData.CreatedAt,
	}

	h.logger.Info("successfully created a new message")
	// if h.messageBroadcaster != nil {
	err = h.BroadcastMessageEvent(ctx, "CREATE_MESSAGE", data)
	if err != nil {
		h.logger.Error("Failed to broadcast tournament match event: ", err)
	}
	// }

	h.logger.Debug("Successfully broad cast message")
	// ctx.JSON(http.StatusAccepted, data)
}

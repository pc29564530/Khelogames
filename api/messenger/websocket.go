package messenger

import (
	"context"
	"encoding/json"
	"fmt"
	db "khelogames/database"

	"khelogames/core/token"
	"khelogames/pkg"
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

func (s *MessageServer) StartMessageHub() {
	for {
		select {
		case message := <-s.messageBroadCast:
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

func (s *MessageServer) StartCricketScoreHub() {
	for {
		select {
		case cricketScore := <-s.scoreBroadCast:
			s.mutex.Lock()
			for client := range s.clients {
				err := client.WriteMessage(websocket.TextMessage, cricketScore)
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
		fmt.Println("Message Type: ", message["type"])
		switch message["type"] {
		case "CREATE_MESSAGE":
			getMessageHub(h, ctx, msg, message["payload"].(map[string]interface{}))
		case "UPDATE_SCORE":
			getCricketScoreHub(h, ctx, msg, message)
		case "INNING_STATUS":
			getCricketScoreHub(h, ctx, msg, message)
		}

		tx, err := h.store.BeginTx(ctx)
		if err != nil {
			h.logger.Error("Failed to begin transactions: ", err)
			return
		}

		defer tx.Rollback()

		h.logger.Debug("Successfully broad cast message")

		err = tx.Commit()
		if err != nil {
			h.logger.Error("Failed to commit the transactions: ", err)
			return
		}
	}
}

func (s *MessageServer) StartRabbitMQConsumer(queueName string) {
	msgs, err := s.rabbitChan.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.logger.Fatal("Failed to register consumer: ", err)
		return
	}
	fmt.Println("QueueName: ", queueName)

	s.logger.Infof("Starting RabbitMQ consumer for queue: %s", queueName)

	go func() {
		for msg := range msgs {
			var message map[string]interface{}
			err := json.Unmarshal(msg.Body, &message)
			if err != nil {
				s.logger.Error("Failed to unmarshal RabbitMQ message: ", err)
				continue
			}

			messageType, ok := message["type"].(string)
			if !ok {
				s.logger.Error("Message type not found or invalid")
				continue
			}

			switch queueName {
			case "chatHub":
				if messageType == "CREATE_MESSAGE" {
					s.messageBroadCast <- msg.Body
					s.logger.Debug("Broadcasted chat message")
				}
			case "scoreHub":
				switch messageType {
				case "ADD_BATSMAN", "ADD_BOWLER", "INNING_STATUS", "UPDATE_SCORE":
					s.scoreBroadCast <- msg.Body
					s.logger.Info("Broadcast score update: ", msg.Body)
				default:
					s.logger.Warn("Unknown score message type: %s", messageType)
				}
			default:
				s.logger.Warnf("Unknown queue: %s", queueName)
			}
		}
	}()
}

func (s *MessageServer) BroadcastFootballEvent(ctx *gin.Context, eventType string, payload map[string]interface{}) error {
	content := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}

	//Log before marshalling
	s.logger.Infof("[BroadcastFootballEvent] Preparing broadcast for eventType=%s", eventType)
	s.logger.Debugf("[BroadcastFootballEvent] Raw payload: %#v", payload)

	body, err := json.Marshal(content)
	if err != nil {
		s.logger.Errorf("failed to marshal message: %v", err)
		return err
	}

	//Log size and body preview
	s.logger.Infof("[BroadcastFootballEvent] Marshaled JSON size: %d bytes", len(body))
	s.logger.Debugf("[BroadcastFootballEvent] Marshaled JSON: %s", string(body))

	//Verify JSON validity before send
	var check map[string]interface{}
	if err := json.Unmarshal(body, &check); err != nil {
		s.logger.Errorf("[BroadcastFootballEvent] Invalid JSON generated: %v", err)
		return err
	}

	//Non-empty check
	if len(body) == 0 {
		s.logger.Warn("[BroadcastFootballEvent] Skipping empty broadcast body")
		return fmt.Errorf("Error empty body")
	}

	//Send to channel
	select {
	case s.scoreBroadCast <- body:
		s.logger.Infof("[BroadcastFootballEvent] Sent to scoreBroadCast successfully (len=%d)", len(s.scoreBroadCast))
	default:
		s.logger.Warn("[BroadcastFootballEvent] scoreBroadCast channel is full or blocked — message dropped")
	}
	return nil
}

func (s *MessageServer) BroadcastCricketEvent(ctx *gin.Context, eventType string, payload map[string]interface{}) error {
	content := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}

	//Log before marshalling
	s.logger.Infof("[BroadcastCricketEvent] Preparing broadcast for eventType=%s", eventType)
	s.logger.Debugf("[BroadcastCricketEvent] Raw payload: %#v", payload)

	body, err := json.Marshal(content)
	if err != nil {
		s.logger.Errorf("failed to marshal message: %v", err)
		return err
	}

	//Log size and body preview
	s.logger.Infof("[BroadcastCricketEvent] Marshaled JSON size: %d bytes", len(body))
	s.logger.Debugf("[BroadcastCricketEvent] Marshaled JSON: %s", string(body))

	//Verify JSON validity before send
	var check map[string]interface{}
	if err := json.Unmarshal(body, &check); err != nil {
		s.logger.Errorf("[BroadcastCricketEvent] Invalid JSON generated: %v", err)
		return err
	}

	//Non-empty check
	if len(body) == 0 {
		s.logger.Warn("[BroadcastCricketEvent] Skipping empty broadcast body")
		return fmt.Errorf("Error empty body")
	}

	//Send to channel
	select {
	case s.scoreBroadCast <- body:
		s.logger.Infof("[BroadcastCricketEvent] Sent to scoreBroadCast successfully (len=%d)", len(s.scoreBroadCast))
	default:
		s.logger.Warn("[BroadcastCricketEvent] scoreBroadCast channel is full or blocked — message dropped")
	}
	return nil
}

func (s *MessageServer) BroadcastTournamentEvent(ctx *gin.Context, eventType string, payload map[string]interface{}) error {
	content := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}

	//Log before marshalling
	s.logger.Infof("[BroadcastTournamentEvent] Preparing broadcast for eventType=%s", eventType)
	s.logger.Debugf("[BroadcastTournamentEvent] Raw payload: %#v", payload)

	body, err := json.Marshal(content)
	if err != nil {
		s.logger.Errorf("failed to marshal message: %v", err)
		return err
	}

	//Log size and body preview
	s.logger.Infof("[BroadcastTournamentEvent] Marshaled JSON size: %d bytes", len(body))
	s.logger.Debugf("[BroadcastTournamentEvent] Marshaled JSON: %s", string(body))

	//Verify JSON validity before send
	var check map[string]interface{}
	if err := json.Unmarshal(body, &check); err != nil {
		s.logger.Errorf("[BroadcastTournamentEvent] Invalid JSON generated: %v", err)
		return err
	}

	//Non-empty check
	if len(body) == 0 {
		s.logger.Warn("[BroadcastTournamentEvent] Skipping empty broadcast body")
		return fmt.Errorf("Error empty body")
	}

	//Send to channel
	select {
	case s.scoreBroadCast <- body:
		s.logger.Infof("[BroadcastTournamentEvent] Sent to scoreBroadCast successfully (len=%d)", len(s.scoreBroadCast))
	default:
		s.logger.Warn("[BroadcastTournamentEvent] scoreBroadCast channel is full or blocked — message dropped")
	}
	return nil
}

func getMessageHub(h *MessageServer, ctx context.Context, msg []byte, message map[string]interface{}) {
	err := h.rabbitChan.PublishWithContext(
		ctx,
		"",
		"chatHub",
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
	var routeContext *gin.Context
	authToken := routeContext.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
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

	arg := db.CreateNewMessageParams{
		SenderID:   authToken.PublicID,
		ReceiverID: receiverPublicID,
		Content:    message["content"].(string),
		MediaUrl:   message["media_url"].(string),
		MediaType:  message["media_type"].(string),
		SentAt:     sentAt,
	}

	h.logger.Debug("create new message params: ", arg)

	_, err = h.store.CreateNewMessage(ctx, arg)
	if err != nil {
		h.logger.Error("unable to store new message: ", err)
		return
	}

	h.logger.Info("successfully created a new message")

	h.messageBroadCast <- msg

	h.logger.Debug("Successfully broad cast message")

}

func getCricketScoreHub(h *MessageServer, ctx *gin.Context, msg []byte, message map[string]interface{}) {
	err := h.rabbitChan.PublishWithContext(
		ctx,
		"",
		"scoreHub",
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

	if h.cricketUpdater == nil {
		h.logger.Error("cricket updater not initialized")
		return
	}

	var cricketData map[string]interface{}
	var inningStatus string

	switch message["payload"].(map[string]interface{})["event_type"].(string) {
	case "normal":
		cricketData, inningStatus = h.cricketUpdater.UpdateInningScoreWS(ctx, message["payload"].(map[string]interface{}))
	case "wide":
		cricketData = h.cricketUpdater.UpdateWideBallWS(ctx, message["payload"].(map[string]interface{}))
	case "no_ball":
		cricketData = h.cricketUpdater.UpdateNoBallsRunsWS(ctx, message["payload"].(map[string]interface{}))
	case "wicket":
		cricketData = h.cricketUpdater.AddCricketWicketsWS(ctx, message["payload"].(map[string]interface{}))
	}

	fmt.Println("Cricket Data: ", cricketData)

	h.logger.Info("successfully created a update score")

	scoreByte, _ := json.Marshal(cricketData)

	h.scoreBroadCast <- scoreByte
	fmt.Println("Ining Stauts: ", inningStatus)

	if inningStatus == "completed" {
		data := h.cricketUpdater.UpdateCricketInningStatusWS(ctx, message["payload"].(map[string]interface{}))
		fmt.Println("Inning Status Completed Dta: ", data)
		var inningStatusByte []byte
		inningStatusByte, err := json.Marshal(data)
		if err != nil {
			h.logger.Error("Failed to marshal :", err)
			return
		}
		h.scoreBroadCast <- inningStatusByte
	}

}

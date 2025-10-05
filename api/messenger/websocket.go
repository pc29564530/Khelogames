package messenger

import (
	"encoding/json"
	"fmt"
	cricket "khelogames/api/sports/cricket"
	db "khelogames/database"
	"khelogames/database/models"

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
			h.logger.Error("Failed to begin transcation: ", err)
			return
		}

		defer tx.Rollback()

		h.logger.Debug("Successfully broad cast message")

		err = tx.Commit()
		if err != nil {
			h.logger.Error("Failed to commit the transcation: ", err)
			return
		}
	}
}

func getMessageHub(h *MessageServer, ctx *gin.Context, msg []byte, message map[string]interface{}) {
	err := h.rabbitChan.PublishWithContext(
		ctx,
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

	fmt.Println("Line no 246: ")

	cricketServer := cricket.NewCricketServer(h.store, h.logger)

	var cricketData map[string]interface{}
	var inningStatus string

	switch message["payload"].(map[string]interface{})["event_type"].(string) {
	case "normal":
		cricketData, inningStatus = cricketServer.UpdateInningScoreWS(ctx, message["payload"].(map[string]interface{}))
	case "wide":
		cricketData = cricketServer.UpdateWideBallWS(ctx, message["payload"].(map[string]interface{}))
	case "no_ball":
		cricketData = cricketServer.UpdateNoBallsRunsWS(ctx, message["payload"].(map[string]interface{}))
	case "wicket":
		cricketData = cricketServer.AddCricketWicketsWS(ctx, message["payload"].(map[string]interface{}))
	}

	// cricketData := cricketServer.UpdateInningScoreWS(ctx, message["payload"].(map[string]interface{}))

	fmt.Println("Cricket Data: ", cricketData)

	h.logger.Info("successfully created a update score")

	scoreByte, _ := json.Marshal(cricketData)

	h.scoreBroadCast <- scoreByte

	fmt.Println("Inning Status 274: ", inningStatus)

	h.logger.Debug("Successfully broad cast score")

	if inningStatus == "completed" {
		fmt.Println("Message: ", message)
		content := message["payload"].(map[string]interface{})
		matchPublicID, err := uuid.Parse(content["match_public_id"].(string))
		if err != nil {
			h.logger.Error("Invalid match UUID format", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match UUID format"})
			return
		}

		batsmanTeamPublicID, err := uuid.Parse(content["batsman_team_public_id"].(string))
		if err != nil {
			h.logger.Error("Invalid batsman team UUID format", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batsman team UUID format"})
			return
		}

		fmt.Println("Message Data: ", content["bowler_public_id"].(string))
		bowlerPublicID, err := uuid.Parse(content["bowler_public_id"].(string))
		if err != nil {
			h.logger.Error("Invalid bowler UUID format", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bowler UUID format"})
			return
		}
		inningNumber := int(content["inning_number"].(float64))

		currBatsman, err := h.store.GetCurrentBattingBatsman(ctx, matchPublicID, batsmanTeamPublicID, inningNumber)
		if err != nil {
			h.logger.Error("Faield to end the inning ")
		}

		var strikerResponse models.BatsmanScore
		var nonStrikerResponse models.BatsmanScore

		for _, curr := range currBatsman {
			if curr.IsStriker {
				strikerResponse = curr
			} else {
				nonStrikerResponse = curr
			}
		}
		fmt.Println("Lien no 309")

		inningScore, _, bowlerResponse, err := h.store.UpdateInningEndStatusByPublicID(ctx, matchPublicID, batsmanTeamPublicID, inningNumber)
		if err != nil {
			h.logger.Error("Faield to end the inning ")
		}
		fmt.Println("INNING SCORE: ", inningScore)
		fmt.Println("BowlerResoponse; ", bowlerResponse)

		strikerPlayerData, err := h.store.GetPlayerByID(ctx, int64(strikerResponse.BatsmanID))
		if err != nil {
			h.logger.Error("Failed to get striker player: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get striker data"})
			return
		}

		nonStrikerPlayerData, err := h.store.GetPlayerByID(ctx, int64(nonStrikerResponse.BatsmanID))
		if err != nil {
			h.logger.Error("Failed to get non-striker player: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get non-striker data"})
			return
		}

		bowlerPlayerData, err := h.store.GetPlayerByPublicID(ctx, bowlerPublicID)
		if err != nil {
			h.logger.Error("Failed to get bowler player: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bowler data"})
			return
		}

		// Build response objects
		striker := map[string]interface{}{
			"player": map[string]interface{}{
				"id":        strikerPlayerData.ID,
				"public_id": strikerPlayerData.PublicID,
				"name":      strikerPlayerData.Name,
				"slug":      strikerPlayerData.Slug,
				"shortName": strikerPlayerData.ShortName,
				"position":  strikerPlayerData.Positions,
			},
			"id":                   strikerResponse.ID,
			"public_id":            strikerResponse.PublicID,
			"match_id":             strikerResponse.MatchID,
			"team_id":              strikerResponse.TeamID,
			"batsman_id":           strikerResponse.BatsmanID,
			"runs_scored":          strikerResponse.RunsScored,
			"balls_faced":          strikerResponse.BallsFaced,
			"fours":                strikerResponse.Fours,
			"sixes":                strikerResponse.Sixes,
			"batting_status":       strikerResponse.BattingStatus,
			"is_striker":           strikerResponse.IsStriker,
			"is_currently_batting": strikerResponse.IsCurrentlyBatting,
			"inning_number":        strikerResponse.InningNumber,
		}

		nonStriker := map[string]interface{}{
			"player": map[string]interface{}{
				"id":        nonStrikerPlayerData.ID,
				"public_id": nonStrikerPlayerData.PublicID,
				"name":      nonStrikerPlayerData.Name,
				"slug":      nonStrikerPlayerData.Slug,
				"shortName": nonStrikerPlayerData.ShortName,
				"position":  nonStrikerPlayerData.Positions,
			},
			"id":                   nonStrikerResponse.ID,
			"public_id":            nonStrikerResponse.PublicID,
			"match_id":             nonStrikerResponse.MatchID,
			"team_id":              nonStrikerResponse.TeamID,
			"batsman_id":           nonStrikerResponse.BatsmanID,
			"runs_scored":          nonStrikerResponse.RunsScored,
			"balls_faced":          nonStrikerResponse.BallsFaced,
			"fours":                nonStrikerResponse.Fours,
			"sixes":                nonStrikerResponse.Sixes,
			"batting_status":       nonStrikerResponse.BattingStatus,
			"is_striker":           nonStrikerResponse.IsStriker,
			"is_currently_batting": nonStrikerResponse.IsCurrentlyBatting,
			"inning_number":        nonStrikerResponse.InningNumber,
		}

		bowler := map[string]interface{}{
			"player": map[string]interface{}{
				"id":        bowlerPlayerData.ID,
				"public_id": bowlerPlayerData.PublicID,
				"name":      bowlerPlayerData.Name,
				"slug":      bowlerPlayerData.Slug,
				"shortName": bowlerPlayerData.ShortName,
				"position":  bowlerPlayerData.Positions,
			},
			"id":                bowlerResponse.ID,
			"public_id":         bowlerResponse.PublicID,
			"match_id":          bowlerResponse.MatchID,
			"team_id":           bowlerResponse.TeamID,
			"bowler_id":         bowlerResponse.BowlerID,
			"ball_number":       bowlerResponse.BallNumber,
			"runs":              bowlerResponse.Runs,
			"wide":              bowlerResponse.Wide,
			"no_ball":           bowlerResponse.NoBall,
			"wickets":           bowlerResponse.Wickets,
			"bowling_status":    bowlerResponse.BowlingStatus,
			"is_current_bowler": bowlerResponse.IsCurrentBowler,
			"inning_number":     bowlerResponse.InningNumber,
		}
		inningPayload := map[string]interface{}{
			"id":                  inningScore.ID,
			"public_id":           inningScore.PublicID,
			"match_id":            inningScore.MatchID,
			"team_id":             inningScore.TeamID,
			"inning_number":       inningScore.InningNumber,
			"score":               inningScore.Score,
			"wickets":             inningScore.Wickets,
			"overs":               inningScore.Overs,
			"run_rate":            inningScore.RunRate,
			"target_run_rate":     inningScore.TargetRunRate,
			"follow_on":           inningScore.FollowOn,
			"is_inning_completed": inningScore.IsInningCompleted,
			"declared":            inningScore.Declared,
		}

		data := map[string]interface{}{
			"type": "INNING_STATUS",
			"payload": map[string]interface{}{
				"striker":       striker,
				"non_striker":   nonStriker,
				"bowler":        bowler,
				"inning_score":  inningPayload,
				"inning_status": "completed",
			},
		}

		fmt.Println("Line no 446: Status; ", inningStatus)

		inningStatusByte, err := json.Marshal(data)
		if err != nil {
			h.logger.Error("Failed to marshal :", err)
			return
		}

		h.scoreBroadCast <- inningStatusByte

	}

}

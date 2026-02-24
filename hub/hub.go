package hub

import (
	"encoding/json"
	"fmt"
	"khelogames/util"

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

func (s *Hub) StartMessageHub() {
	s.logger.Info("StartMessageHub started")
	for {
		select {
		case message := <-s.MessageBroadcast:
			var data map[string]interface{}
			err := json.Unmarshal(message, &data)
			if err != nil {
				s.logger.Error("Failed to unmarshal message:", err)
				continue
			}

			payload, ok := data["payload"].(map[string]interface{})
			if !ok {
				s.logger.Error("invalid payload structure")
				continue
			}

			receiverID, _ := payload["receiver"].(map[string]interface{})["public_id"].(string)
			senderID, _ := payload["sender"].(map[string]interface{})["public_id"].(string)
			s.logger.Debugf("Broadcasting message: senderID=%s receiverID=%s", senderID, receiverID)

			// Copy subscriber sets under lock to avoid race conditions
			s.mu.Lock()
			senSubsCopy := make([]*Client, 0, len(s.subscriber[senderID]))
			for c := range s.subscriber[senderID] {
				senSubsCopy = append(senSubsCopy, c)
			}
			recSubsCopy := make([]*Client, 0, len(s.subscriber[receiverID]))
			for c := range s.subscriber[receiverID] {
				recSubsCopy = append(recSubsCopy, c)
			}
			s.mu.Unlock()

			var failedClients []*Client
			for _, client := range senSubsCopy {
				if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					failedClients = append(failedClients, client)
				}
			}
			for _, client := range recSubsCopy {
				if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					failedClients = append(failedClients, client)
				}
			}
			if len(failedClients) > 0 {
				s.mu.Lock()
				for _, client := range failedClients {
					delete(s.Clients, client)
					client.Conn.Close()
				}
				s.mu.Unlock()
			}
		}
	}
}

func (s *Hub) StartTournamentHub() {

	s.logger.Info("StartTournamentHub started")
	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("StartTournamentHub panic: %v", r)
		}
	}()
	for {
		select {
		case tournament := <-s.TournamentBroadcast:
			s.mu.Lock()
			for client := range s.Clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, tournament)
				if err != nil {
					s.logger.Errorf("Failed to write tournament message to client: %v", err)
					delete(s.Clients, client)
					client.Conn.Close()
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Hub) StartCricketHub() {

	s.logger.Info("StartCricketHub started")
	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("StartCricketHub panic: %v", r)
		}
	}()
	for {
		select {
		case cricket := <-s.CricketBroadcast:
			s.mu.Lock()
			for client := range s.Clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, cricket)
				if err != nil {
					s.logger.Errorf("Failed to write cricket message to client: %v", err)
					delete(s.Clients, client)
					client.Conn.Close()
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Hub) StartFootballHub() {
	s.logger.Info("StartFootballHub started")
	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("StartFootballHub panic: %v", r)
		}
	}()
	for {
		select {
		case football := <-s.FootballBroadcast:
			s.mu.Lock()
			for client := range s.Clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, football)
				if err != nil {
					s.logger.Errorf("Failed to write football message to client: %v", err)
					delete(s.Clients, client)
					client.Conn.Close()
				}
			}
			s.mu.Unlock()
		}
	}
}

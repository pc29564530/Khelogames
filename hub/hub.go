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
	fmt.Println("Start message hub ")
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
			fmt.Println("REciever ID: ", receiverID)
			fmt.Println("Sender ID: ", senderID)
			s.mu.Lock()
			senSubs := s.subscriber[senderID]
			recSubs := s.subscriber[receiverID]
			s.mu.Unlock()

			fmt.Println("Sender Sub: ", senSubs)
			fmt.Println("Receiver Sub: ", recSubs)

			for client := range senSubs {
				err := client.Conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					delete(s.Clients, client)
					client.Conn.Close()
				}
			}

			for client := range recSubs {
				err := client.Conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					delete(s.Clients, client)
					client.Conn.Close()
				}
			}
		}
	}
}

func (s *Hub) StartTournamentHub() {
	// fmt.Println("Lien noi 41 Start Tournament Hub")
	// fmt.Println("Subscriber: ", s.subscriber)
	// fmt.Println("Cleint: ", s.Clients)
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
				fmt.Println("Client: ", client)
				err := client.Conn.WriteMessage(websocket.TextMessage, tournament)
				if err != nil {
					fmt.Println("Failed to write message to tournament: client:", err)
					delete(s.Clients, client)
					client.Conn.Close()
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Hub) StartCricketHub() {
	// fmt.Println("Lien noi 41 Start Tournament Hub")
	// fmt.Println("Subscriber: ", s.subscriber)
	// fmt.Println("Cleint: ", s.Clients)
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
				fmt.Println("Client: ", client)
				err := client.Conn.WriteMessage(websocket.TextMessage, cricket)
				if err != nil {
					fmt.Println("Failed to write message to cricket: client:", err)
					delete(s.Clients, client)
					client.Conn.Close()
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Hub) StartFootballHub() {
	// fmt.Println("Lien noi 41 Start Tournament Hub")
	// fmt.Println("Subscriber: ", s.subscriber)
	// fmt.Println("Cleint: ", s.Clients)
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
				fmt.Println("Client: ", client)
				err := client.Conn.WriteMessage(websocket.TextMessage, football)
				if err != nil {
					fmt.Println("Failed to write message to cricket: client:", err)
					delete(s.Clients, client)
					client.Conn.Close()
				}
			}
			s.mu.Unlock()
		}
	}
}

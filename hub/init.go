package hub

import (
	"khelogames/api/shared"
	"khelogames/core/token"
	"khelogames/database"
	"khelogames/logger"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	ampq "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	Conn         *websocket.Conn
	UserPublicID uuid.UUID
	SendChan     chan []byte
}

type Hub struct {
	Clients map[*Client]bool
	mu      sync.Mutex

	MessageBroadcast    chan []byte
	CricketBroadcast    chan []byte
	FootballBroadcast   chan []byte
	TournamentBroadcast chan []byte

	logger             *logger.Logger
	store              *database.Store
	upgrader           websocket.Upgrader
	rabbitChan         *ampq.Channel
	tokenMaker         token.Maker
	scoreBroadcaster   *shared.ScoreBroadcaster
	messageBroadcaster shared.MessageBroadcaster
	subscriber         map[string]map[*Client]bool
}

// NewHub creates a new Hub instance and starts broadcaster loops
func NewHub(store *database.Store, logger *logger.Logger, upgrader websocket.Upgrader, rabbitChan *ampq.Channel, tokenMaker token.Maker, scoreBroadcaster shared.ScoreBroadcaster, messageBroadcaster shared.MessageBroadcaster, subscriber map[string]map[*Client]bool) *Hub {
	h := &Hub{
		Clients:             make(map[*Client]bool),
		MessageBroadcast:    make(chan []byte, 256),
		CricketBroadcast:    make(chan []byte),
		FootballBroadcast:   make(chan []byte),
		TournamentBroadcast: make(chan []byte),
		logger:              logger,
		store:               store,
		upgrader:            upgrader,
		rabbitChan:          rabbitChan,
		tokenMaker:          tokenMaker,
		scoreBroadcaster:    &scoreBroadcaster,
		messageBroadcaster:  messageBroadcaster,
		subscriber:          subscriber,
	}

	go h.StartMessageHub()
	go h.StartCricketHub()
	go h.StartFootballHub()
	go h.StartTournamentHub()

	h.logger.Info("Hub initialized successfully")
	return h
}

func (h *Hub) AddClient(conn *websocket.Conn, userPublicID uuid.UUID) *Client {
	client := &Client{
		Conn:         conn,
		UserPublicID: userPublicID,
		SendChan:     make(chan []byte, 256),
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.Clients[client] = true

	h.logger.Infof("Added WebSocket client: %s", client.UserPublicID)
	return client
}

func (h *Hub) RemoveClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.Clients, client)
	// Guard against double-close of the send channel
	select {
	case _, ok := <-client.SendChan:
		if ok {
			close(client.SendChan)
		}
	default:
		close(client.SendChan)
	}
	client.Conn.Close()
	h.logger.Infof("Removed WebSocket client: %s", client.UserPublicID)
}

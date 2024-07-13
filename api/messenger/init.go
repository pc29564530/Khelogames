package messenger

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/token"
	"sync"

	"github.com/gorilla/websocket"
	ampq "github.com/rabbitmq/amqp091-go"
)

type MessageServer struct {
	store      *db.Store
	tokenMaker token.Maker
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	rabbitChan *ampq.Channel
	mutex      sync.Mutex
	logger     *logger.Logger
}

func NewMessageServer(store *db.Store, tokenMaker token.Maker, clients map[*websocket.Conn]bool, broadcast chan []byte, upgrader websocket.Upgrader, rabbitChan *ampq.Channel, logger *logger.Logger) *MessageServer {
	return &MessageServer{
		store:      store,
		tokenMaker: tokenMaker,
		upgrader:   upgrader,
		clients:    clients,
		broadcast:  broadcast,
		rabbitChan: rabbitChan,
		logger:     logger,
	}
}

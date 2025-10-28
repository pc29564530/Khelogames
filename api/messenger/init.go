package messenger

import (
	shared "khelogames/api/shared"
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/logger"
	"sync"

	"github.com/gorilla/websocket"
	ampq "github.com/rabbitmq/amqp091-go"
)

type MessageServer struct {
	store            *db.Store
	tokenMaker       token.Maker
	upgrader         websocket.Upgrader
	clients          map[*websocket.Conn]bool
	messageBroadCast chan []byte
	scoreBroadCast   chan []byte
	rabbitChan       *ampq.Channel
	mutex            sync.Mutex
	logger           *logger.Logger
	cricketUpdater   shared.CricketScoreUpdater
	scoreBroadcaster shared.ScoreBroadcaster
}

func NewMessageServer(store *db.Store, tokenMaker token.Maker, clients map[*websocket.Conn]bool, messageBroadCast chan []byte, scoreBroadCast chan []byte, upgrader websocket.Upgrader, rabbitChan *ampq.Channel, logger *logger.Logger, cricketUpdater shared.CricketScoreUpdater, scoreBroadcaster shared.ScoreBroadcaster) *MessageServer {
	return &MessageServer{
		store:            store,
		tokenMaker:       tokenMaker,
		upgrader:         upgrader,
		clients:          clients,
		messageBroadCast: messageBroadCast,
		scoreBroadCast:   scoreBroadCast,
		rabbitChan:       rabbitChan,
		logger:           logger,
		cricketUpdater:   cricketUpdater,
		scoreBroadcaster: scoreBroadcaster,
	}
}

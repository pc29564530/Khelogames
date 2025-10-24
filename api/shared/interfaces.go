package shared

import (
	"khelogames/logger"

	"github.com/gin-gonic/gin"
	ampq "github.com/rabbitmq/amqp091-go"
)

// ScoreBroadcaster defines the interface for broadcasting cricket score updates
type ScoreBroadcaster interface {
	BroadcastCricketEvent(ctx *gin.Context, eventType string, payload map[string]interface{}) error
	BroadcastFootballEvent(ctx *gin.Context, eventType string, payload map[string]interface{}) error
	BroadcastTournamentEvent(ctx *gin.Context, eventType string, payload map[string]interface{}) error
}

// CricketScoreUpdater defines the interface for updating cricket scores
type CricketScoreUpdater interface {
	UpdateInningScoreWS(ctx *gin.Context, message map[string]interface{}) (data map[string]interface{}, inningStatus string)
	UpdateWideBallWS(ctx *gin.Context, message map[string]interface{}) map[string]interface{}
	UpdateNoBallsRunsWS(ctx *gin.Context, message map[string]interface{}) map[string]interface{}
	AddCricketWicketsWS(ctx *gin.Context, message map[string]interface{}) map[string]interface{}
	UpdateCricketInningStatusWS(ctx *gin.Context, message map[string]interface{}) (data map[string]interface{})
}

// ScoreServiceConfig holds the configuration for score services
type ScoreServiceConfig struct {
	Logger     *logger.Logger
	RabbitChan *ampq.Channel
}

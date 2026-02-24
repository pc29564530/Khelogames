package hub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	ampq "github.com/rabbitmq/amqp091-go"
)

func (s *Hub) BroadcastFootballEvent(ctx *gin.Context, eventType string, payload map[string]interface{}) error {
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
	case s.FootballBroadcast <- body:
		s.logger.Infof("[BroadcastFootballEvent] Sent to scoreBroadCast successfully (len=%d)", len(s.FootballBroadcast))
	default:
		s.logger.Warn("[BroadcastFootballEvent] scoreBroadCast channel is full or blocked — message dropped")
	}
	return nil
}

func (s *Hub) BroadcastCricketEvent(ctx *gin.Context, eventType string, payload map[string]interface{}) error {
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
	case s.CricketBroadcast <- body:
		s.logger.Infof("[BroadcastCricketEvent] Sent to scoreBroadCast successfully (len=%d)", len(s.CricketBroadcast))
	default:
		s.logger.Warn("[BroadcastCricketEvent] scoreBroadCast channel is full or blocked — message dropped")
	}
	return nil
}

func (s *Hub) BroadcastMessageEvent(ctx context.Context, eventType string, payload map[string]interface{}) error {
	content := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}

	//Log before marshalling
	s.logger.Infof("[BroadcastMessageEvent] Preparing broadcast for eventType=%s", eventType)
	s.logger.Debugf("[BroadcastMessageEvent] Raw payload: %#v", payload)

	body, err := json.Marshal(content)
	if err != nil {
		s.logger.Errorf("failed to marshal message: %v", err)
		return err
	}

	if s.rabbitChan == nil {
		s.logger.Warn("[BroadcastMessageEvent] RabbitMQ not available, skipping publish")
		select {
		case s.MessageBroadcast <- body:
			s.logger.Info("[BroadcastMessageEvent] Sent directly to MessageBroadcast channel")
		default:
			s.logger.Warn("[BroadcastMessageEvent] MessageBroadcast channel full — message dropped")
		}
		return nil
	}

	err = s.rabbitChan.PublishWithContext(
		ctx,
		"",
		"chatHub",
		false,
		false,
		ampq.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		s.logger.Errorf("failed to publish RabbitMQ: %v", err)
	}

	//Log size and body preview
	s.logger.Infof("[BroadcastMessageEvent] Marshaled JSON size: %d bytes", len(body))
	s.logger.Debugf("[BroadcastMessageEvent] Marshaled JSON: %s", string(body))

	return nil
}

func (s *Hub) BroadcastTournamentEvent(ctx *gin.Context, eventType string, payload map[string]interface{}) error {

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

	//Non-empty check
	if len(body) == 0 {
		s.logger.Warn("[BroadcastTournamentEvent] Skipping empty broadcast body")
		return fmt.Errorf("Error empty body")
	}

	var check map[string]interface{}
	if err := json.Unmarshal(body, &check); err != nil {
		s.logger.Errorf("[BroadcastFootballEvent] Invalid JSON generated: %v", err)
		return err
	}

	s.logger.Infof("[BroadcastTournamentEvent] Marshaled JSON size: %d bytes", len(body))
	s.logger.Debugf("[BroadcastTournamentEvent] Marshaled JSON: %s", string(body))

	//Send to channel
	select {
	case s.TournamentBroadcast <- body:
		s.logger.Infof("[BroadcastTournamentEvent] Sent to scoreBroadCast successfully (len=%d)", len(s.TournamentBroadcast))
	default:
		s.logger.Warn("[BroadcastTournamentEvent] scoreBroadCast channel is full or blocked — message dropped")
	}
	return nil
}

package hub

import (
	"encoding/json"
)

func (s *Hub) StartRabbitMQConsumer(queueName string) {
	if s.rabbitChan == nil {
		s.logger.Warnf("RabbitMQ not available, skipping consumer for queue: %s", queueName)
		return
	}

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

			s.logger.Debugf("RabbitMQ message received: type=%s queue=%s", messageType, queueName)
			switch queueName {
			case "scoreHub":
				s.logger.Warnf("No Score hub")
			case "chatHub":
				switch messageType {
				case "CREATE_MESSAGE":
					s.MessageBroadcast <- msg.Body
				default:
					s.logger.Warn("Unknown message type: %s", messageType)
				}
			default:
				s.logger.Warnf("Unknown queue: %s", queueName)
			}
		}
	}()
}

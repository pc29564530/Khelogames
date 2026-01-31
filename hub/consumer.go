package hub

import (
	"encoding/json"
	"fmt"
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

			fmt.Println("Unmarshal Message: ", message)

			messageType, ok := message["type"].(string)
			if !ok {
				s.logger.Error("Message type not found or invalid")
				continue
			}

			fmt.Println("Message Type: ", messageType)
			fmt.Println("Queue Name: ", queueName)
			switch queueName {
			// case "chatHub":
			// 	if messageType == "CREATE_MESSAGE" {
			// 		s.MessageBroadcast <- msg.Body
			// 		s.logger.Debug("Broadcasted chat message")
			// // 	}
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

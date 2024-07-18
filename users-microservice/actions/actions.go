package actions

import (
	"encoding/json"
	"log"
	"users/conn"
	"users/models"

	"github.com/IBM/sarama"
)

func PushEventToKafka(topic string, user *models.User) {
	userData, err := json.Marshal(user)
	if err != nil {
		log.Printf("failed to marshal user data: %v", err)
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(userData),
	}

	partition, offset, err := conn.GetKafkaProducer().SendMessage(msg)
	if err != nil {
		log.Printf("failed to produce message to Kafka: %v", err)
	} else {
		log.Printf("message sent to partition %d at offset %d", partition, offset)
	}
}

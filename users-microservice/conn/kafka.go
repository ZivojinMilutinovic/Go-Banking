package conn

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

var kafkaProducer sarama.SyncProducer

func ConnectKafka() {
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 25
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	prd, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKER")}, config)
	if err != nil {
		panic("Kafka connection has not been intitialized:" + err.Error())
	} else {
		log.Println("Kafka conntection sucessfully established")
	}

	kafkaProducer = prd
}

func GetKafkaProducer() sarama.SyncProducer {
	if kafkaProducer == nil {
		panic("Kafka producer has not been Initialized")
	}

	return kafkaProducer
}

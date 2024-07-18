package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transactions/api"
	"transactions/conn"
	"transactions/controllers"
	"transactions/models"

	"github.com/IBM/sarama"
	"github.com/nats-io/nats.go"
)

func main() {
	time.Sleep(20 * time.Second)
	log.Println("Starting transactions microservice")
	conn.ConnectDB()
	conn.AutoMigrate()

	nc, err := nats.Connect(os.Getenv("NATS_SERVER"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	api := api.GetUserApi(nc)
	api.SendUserBalance()

	consumer, err := sarama.NewConsumer([]string{os.Getenv("KAFKA_BROKER")}, nil)
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()
	topic := "user-created-topic"

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("failed to start consumer for topic %s: %v", topic, err)
	}
	defer partitionConsumer.Close()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// Consume messages
	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				var user models.User
				if err := json.Unmarshal(msg.Value, &user); err != nil {
					log.Printf("failed to unmarshal message: %v", err)
				} else {
					log.Printf("Received user: %+v\n", user)
					result, err := api.CreateUser(&user)
					if err != nil {
						log.Printf("User was not created error: %s", err.Error())
					} else {
						log.Printf("User created successfully:%v", result)
					}
				}
			case <-sigterm:
				log.Println("Shutting down consumer...")
				return
			}
		}
	}()

	go func() {
		controllers.SetupServer()
	}()

	// Wait for termination signal

	<-sigterm
}

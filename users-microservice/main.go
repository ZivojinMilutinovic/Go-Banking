package main

import (
	"log"
	"os"
	"time"
	"users/conn"
	"users/controllers"

	"github.com/nats-io/nats.go"
)

func main() {
	time.Sleep(20 * time.Second)
	log.Println("Starting user microservice")
	nc, err := nats.Connect(os.Getenv("NATS_SERVER"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	conn.ConnectDB()
	conn.AutoMigrate()
	conn.ConnectKafka()
	controllers.SetupServer(nc)
}

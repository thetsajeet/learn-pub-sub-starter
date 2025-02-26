package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

const CONNECTION_STRING = "amqp://guest:guest@localhost:5672/"

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	log.Default().Println("Starting Peril Server")

	amqpConnection, err := amqp.Dial(CONNECTION_STRING)
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConnection.Close()

	log.Default().Println("Connection successful to AMQP")

	<-signalChan
	fmt.Println()
	log.Default().Println("Stopping Peril Server")
}

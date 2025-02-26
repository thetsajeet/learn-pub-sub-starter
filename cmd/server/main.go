package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/cmd"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	log.Default().Println("Starting Peril Server")

	amqpConnection, err := amqp.Dial(cmd.CONNECTION_STRING)
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConnection.Close()

	amqlChannel, err := amqpConnection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	pubsub.PublishJSON(amqlChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})

	log.Default().Println("Connection successful to AMQP")

	<-signalChan
	fmt.Println()
	log.Default().Println("Stopping Peril Server")
}

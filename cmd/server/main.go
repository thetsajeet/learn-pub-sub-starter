package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/cmd"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	amqpChannel, err := amqpConnection.Channel()
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Println("Connection successful to RabbitMQ server")

	_, _, err = pubsub.DeclareAndBind(amqpConnection, routing.ExchangePerilTopic, routing.GameLogSlug, "game_logs.*", 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Println("Created a new game queue")

	gamelogic.PrintServerHelp()
	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}

		switch inputs[0] {
		case "pause":
			log.Default().Println("Sending pause message")
			if err := pubsub.PublishJSON(amqpChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			}); err != nil {
				fmt.Println(err)
				continue
			}
		case "resume":
			log.Default().Println("Sending a resume message")
			if err := pubsub.PublishJSON(amqpChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			}); err != nil {
				fmt.Println(err)
				continue
			}
		case "quit":
			log.Default().Println("Exitting")
			return
		default:
			fmt.Println("Cannot understand the command")
		}
	}
}

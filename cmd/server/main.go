package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/cmd"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWriteLog() func(gl routing.GameLog) pubsub.AckType {
	return func(gl routing.GameLog) pubsub.AckType {
		defer fmt.Print("> ")
		if err := gamelogic.WriteLog(gl); err != nil {
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}

func main() {
	log.Default().Println("Starting Peril Server")

	amqpConnection, err := amqp.Dial(cmd.CONNECTION_STRING)
	if err != nil {
		fmt.Print('s')
		log.Fatal(err)
	}
	defer amqpConnection.Close()

	amqpChannel, err := amqpConnection.Channel()
	if err != nil {
		fmt.Print('t')
		log.Fatal(err)
	}
	log.Default().Println("Connection successful to RabbitMQ server")

	if err := pubsub.SubscribeGob(
		amqpConnection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		0,
		handlerWriteLog(),
	); err != nil {
		fmt.Print('d')
		log.Fatal(err)
	}

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

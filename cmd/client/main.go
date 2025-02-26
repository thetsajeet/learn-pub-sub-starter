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

	fmt.Println("Starting Peril client...")

	amqpConnection, err := amqp.Dial(cmd.CONNECTION_STRING)
	if err != nil {
		log.Fatal(err)
	}

	welcome, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}

	pubsub.DeclareAndBind(amqpConnection, routing.ExchangePerilDirect, "pause."+welcome, routing.PauseKey, 1)

	<-signalChan
	fmt.Println()
	log.Default().Println("Closing Peril Client...")
}

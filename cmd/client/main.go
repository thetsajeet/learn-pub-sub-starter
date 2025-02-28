package main

import (
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/cmd"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Default().Println("Starting Peril client...")

	amqpConnection, err := amqp.Dial(cmd.CONNECTION_STRING)
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConnection.Close()
	log.Default().Println("Game client connected to RabbitMQ Server")

	publishCh, err := amqpConnection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}

	gs := gamelogic.NewGameState(username)

	if err := pubsub.SubscribeJSON(
		amqpConnection,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		1,
		handlerPause(gs),
	); err != nil {
		log.Fatal(err)
	}
	if err := pubsub.SubscribeJSON(
		amqpConnection,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gs.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		1,
		handlerMove(gs, publishCh),
	); err != nil {
		log.Fatal(err)
	}
	if err = pubsub.SubscribeJSON(
		amqpConnection,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		0,
		handlerWar(gs),
	); err != nil {
		log.Fatal(err)
	}
	clientRepl(gs, publishCh)
}

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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}

	pubsub.DeclareAndBind(amqpConnection, routing.ExchangePerilDirect, "pause."+username, routing.PauseKey, 1)

	gameState := gamelogic.NewGameState(username)

outerloop:
	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}

		switch inputs[0] {
		case "spawn":
			if err := gameState.CommandSpawn(inputs); err != nil {
				fmt.Println(err)
			}
		case "move":
			gameMove, err := gameState.CommandMove(inputs)
			if err != nil {
				fmt.Print(err)
			}
			fmt.Print(gameMove)
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			break outerloop
		default:
			fmt.Println("unknown command")
		}
	}

	<-signalChan
	fmt.Println()
	log.Default().Println("Closing Peril Client...")
}

package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func clientRepl(gs *gamelogic.GameState, publishCh *amqp.Channel) {
	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}

		switch inputs[0] {
		case "spawn":
			if err := gs.CommandSpawn(inputs); err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			gameMove, err := gs.CommandMove(inputs)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if err := pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+gameMove.Player.Username,
				gameMove,
			); err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(inputs) < 2 {
				fmt.Println("require 2 or more arguments")
				continue
			}
			n, err := strconv.Atoi(inputs[1])
			if err != nil {
				fmt.Println("second argument must be a number")
				continue
			}
			for _ = range n {
				mLog := gamelogic.GetMaliciousLog()
				if err := pubsub.PublishGob(
					publishCh,
					routing.ExchangePerilTopic,
					routing.GameLogSlug+"."+gs.GetUsername(),
					routing.GameLog{
						CurrentTime: time.Now(),
						Message:     mLog,
						Username:    gs.GetUsername(),
					},
				); err != nil {
					fmt.Println(err)
				}
			}
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}
}

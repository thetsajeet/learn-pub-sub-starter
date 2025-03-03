package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(mv gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(mv)
		if outcome == gamelogic.MoveOutComeSafe {
			return pubsub.Ack
		} else if outcome == gamelogic.MoveOutcomeMakeWar {
			if err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%s.%s", routing.WarRecognitionsPrefix, gs.GetUsername()),
				gamelogic.RecognitionOfWar{
					Attacker: mv.Player,
					Defender: gs.GetPlayerSnap(),
				},
			); err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		} else {
			return pubsub.NackDiscard
		}
	}
}

func handlerWar(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, _, _ := gs.HandleWar(rw)

		if outcome == gamelogic.WarOutcomeNotInvolved {
			return pubsub.NackRequeue
		} else if outcome == gamelogic.WarOutcomeNoUnits {
			return pubsub.NackDiscard
		} else if outcome == gamelogic.WarOutcomeOpponentWon || outcome == gamelogic.WarOutcomeYouWon || outcome == gamelogic.WarOutcomeDraw {
			msg := ""
			if outcome == gamelogic.WarOutcomeYouWon {
				msg = fmt.Sprintf("%s won a war against %s", rw.Attacker.Username, rw.Defender.Username)
			} else if outcome == gamelogic.WarOutcomeOpponentWon {
				msg = fmt.Sprintf("%s won a war against %s", rw.Defender.Username, rw.Attacker.Username)
			} else {
				msg = fmt.Sprintf("A war between %s and %s resulted in a draw", rw.Attacker.Username, rw.Defender.Username)
			}

			if err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+gs.GetUsername(), routing.GameLog{
				CurrentTime: time.Now(),
				Message:     msg,
				Username:    gs.GetUsername(),
			}); err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		} else {
			fmt.Print("improper outcome")
			return pubsub.NackDiscard
		}
	}
}

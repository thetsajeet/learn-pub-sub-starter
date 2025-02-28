package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](
	ch *amqp.Channel,
	exchange,
	key string,
	val T,
) error {
	jsonVal, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonVal,
	})

	return err
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int,
	table amqp.Table,
) (*amqp.Channel, amqp.Queue, error) {
	amqpChannel, err := conn.Channel()
	if err != nil {
		return amqpChannel, amqp.Queue{}, err
	}

	queue, err := amqpChannel.QueueDeclare(
		queueName,
		simpleQueueType == 0,
		simpleQueueType == 1,
		simpleQueueType == 1,
		false,
		table,
	)
	if err != nil {
		return amqpChannel, queue, nil
	}

	err = amqpChannel.QueueBind(queueName, key, exchange, false, nil)

	return amqpChannel, queue, err
}

type AckType string

const (
	Ack         AckType = "Ack"
	NackRequeue AckType = "NackRequeue"
	NackDiscard AckType = "NackDiscard"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int,
	handler func(T) AckType,
) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType, amqp.Table{
		"x-dead-letter-exchange": "peril_dlx",
	})
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	unmarshaller := func(data []byte) (T, error) {
		var target T
		err := json.Unmarshal(data, &target)
		return target, err
	}

	go func() {
		defer ch.Close()
		for msg := range msgs {
			target, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Println(err)
				continue
			}
			switch ack := handler(target); ack {
			case Ack:
				fmt.Println("ack")
				msg.Ack(false)
			case NackRequeue:
				fmt.Println("nack")
				msg.Nack(false, true)
			case NackDiscard:
				fmt.Println("nack disc")
				msg.Nack(false, false)
			default:
				fmt.Println("unknown ack")
			}
		}
	}()

	return nil
}

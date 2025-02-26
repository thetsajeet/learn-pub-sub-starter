package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
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

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, simpleQueueType int) (*amqp.Channel, amqp.Queue, error) {
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
		nil,
	)
	if err != nil {
		return amqpChannel, queue, nil
	}

	err = amqpChannel.QueueBind(queueName, key, exchange, false, nil)

	return amqpChannel, queue, err
}

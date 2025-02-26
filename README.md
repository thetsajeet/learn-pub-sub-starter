# learn-pub-sub-starter (Peril)

## Notes

- Pub/Sub is a messaging pattern where senders of messages (publishers) do not send messages directly to receivers (subscribers). Instead, they just publish to a single broker. The publisher doesn't need to worry about who all the subscribers are. The broker is responsible for delivering a copy of the message to any interested subscribers.
- Pub/Sub systems are often used to enable "event-driven design", or "event-driven architecture". An event-driven architecture uses events to trigger and communicate between decoupled systems

### RabbitMQ

- RabbitMQ is a popular open-source message broker that implements the AMQP protocol and I've personally been using it happily (well, mostly happily) for years. It's great because it's open-source, flexible, powerful, and (reasonably) easy to use.
- The RabbitMQ server: The message broker itself. It's a server you can run locally (or on a server in production) to send and receive messages.
- Client library: We'll be programming in Go, so our code will use the AMQP library to interact with the RabbitMQ server.
- Management UI: RabbitMQ comes with a management UI that you can use to monitor and manage your RabbitMQ server. It's a web-based UI that you can access in your browser.

### Components of RabbitMQ

- Exchange: A routing agent that sends messages to queues.
- Binding: A link between an exchange and a queue that uses a routing key to decide which messages go to the queue.
- Queue: A buffer in the RabbitMQ server that holds messages until they are consumed.
- Channel: A virtual connection inside a connection that allows you to create queues, exchanges, and publish messages.
- Connection: A TCP connection to the RabbitMQ server.

### Types of Exchanges

- Direct: Messages are routed to the queues based on the message routing key exactly matching the binding key of the queue.
- Topic: Messages are routed to queues based on wildcard matches between the routing key and the routing pattern specified in the binding.
- Fanout: It routes messages to all of the queues bound to it, ignoring the routing key.
- Headers: Routes based on header values instead of the routing key. It's similar to topic but uses message header attributes for routing.

## Queues

- Queues can be "durable" or "transient". Durable queues survive a RabbitMQ server restart, while transient queues do not.
- The metadata of a durable queue is stored on disk, while transient queues are only stored in memory.
- Exclusive: The queue can only be used by the connection that created it.
- Auto-delete: The queue will be automatically deleted when its last connection is closed.

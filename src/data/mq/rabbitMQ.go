package mq

import (
    "fmt"
    "log"
    "user-service/config"

    "github.com/streadway/amqp"
)

var rabbitMQConn *amqp.Connection
var rabbitMQChannels map[string]*amqp.Channel

// InitRabbitMQ initializes the RabbitMQ connection and channels
func InitRabbitMQ(cfg *config.Config) error {
    // Create a connection
    conn, err := amqp.Dial(cfg.RabbitMQ.URL)
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    rabbitMQConn = conn
    rabbitMQChannels = make(map[string]*amqp.Channel)

    // Declare all channels based on the configuration
    for _, queueName := range cfg.RabbitMQ.Queues {
        ch, err := conn.Channel()
        if err != nil {
            conn.Close()
            return fmt.Errorf("failed to open channel for queue '%s': %w", queueName, err)
        }

        _, err = ch.QueueDeclare(
            queueName, // queue name
            true,      // durable
            false,     // auto-deleted when unused
            false,     // not exclusive
            false,     // no-wait
            nil,       // arguments
        )
        if err != nil {
            ch.Close()
            conn.Close()
            return fmt.Errorf("failed to declare queue '%s': %w", queueName, err)
        }

        rabbitMQChannels[queueName] = ch
        log.Printf("RabbitMQ channel initialized for queue: %s", queueName)
    }

    return nil
}

// GetRabbitMQChannel retrieves the channel for a specific queue
func GetRabbitMQChannel(queueName string) (*amqp.Channel, error) {
    ch, ok := rabbitMQChannels[queueName]
    if !ok {
        return nil, fmt.Errorf("queue '%s' not found", queueName)
    }
    return ch, nil
}

// CloseRabbitMQ closes all channels and the RabbitMQ connection
func CloseRabbitMQ() {
    for _, ch := range rabbitMQChannels {
        if err := ch.Close(); err != nil {
            log.Printf("failed to close channel: %v", err)
        }
    }
    if err := rabbitMQConn.Close(); err != nil {
        log.Printf("failed to close RabbitMQ connection: %v", err)
    }
}

// PublishMessage publishes a message to a specified queue
func PublishMessage(queueName string, message []byte) error {
    ch, err := GetRabbitMQChannel(queueName)
    if err != nil {
        return err
    }

    err = ch.Publish(
        "",         // exchange
        queueName,  // routing key (queue name)
        false,      // mandatory
        false,      // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        message,
        },
    )
    if err != nil {
        return fmt.Errorf("failed to publish message to queue '%s': %w", queueName, err)
    }

    log.Printf("Message published to queue: %s", queueName)
    return nil
}

// ConsumeMessages sets up a consumer for a specified queue
func ConsumeMessages(queueName string, handler func([]byte)) error {
    ch, err := GetRabbitMQChannel(queueName)
    if err != nil {
        return err
    }

    msgs, err := ch.Consume(
        queueName, // queue name
        "",        // consumer name
        true,      // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if err != nil {
        return fmt.Errorf("failed to register consumer for queue '%s': %w", queueName, err)
    }

    go func() {
        for d := range msgs {
            handler(d.Body)
        }
    }()

    log.Printf("Consuming messages from queue: %s", queueName)
    return nil
}

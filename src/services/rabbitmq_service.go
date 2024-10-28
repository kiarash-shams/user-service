package services

import (
    "encoding/json"
    "fmt"
	"user-service/pkg/logging"
    "github.com/streadway/amqp"
  
)

// RabbitMQService struct to hold connection and channels
type RabbitMQService struct {
    Connection *amqp.Connection
    Channels   map[string]*amqp.Channel
    logger     logging.Logger // Added logger for logging purposes
}

// NewRabbitMQService initializes the RabbitMQ connection and declares fixed channels
func NewRabbitMQService(url string, logger logging.Logger) (*RabbitMQService, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        logger.Error(logging.RabbitMQ, logging.ExternalService, fmt.Sprintf("Failed to connect to RabbitMQ: %v", err), nil)
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    service := &RabbitMQService{
        Connection: conn,
        Channels:   make(map[string]*amqp.Channel),
        logger:     logger, // Use the provided logger
    }

    // Define the queues we need
    queueNames := []string{"sms_queue", "email_queue", "webhook_queue"}

    // Declare each queue and create a channel for it
    for _, queueName := range queueNames {
        ch, err := conn.Channel()
        if err != nil {
            conn.Close()
            logger.Error(logging.RabbitMQ, logging.ExternalService, fmt.Sprintf("Failed to open a channel for queue '%s': %v", queueName, err), nil)
            return nil, fmt.Errorf("failed to open a channel for queue '%s': %w", queueName, err)
        }

        _, err = ch.QueueDeclare(
            queueName, // queue name
            true,      // durable
            false,     // auto-delete when unused
            false,     // exclusive
            false,     // no-wait
            nil,       // arguments
        )
        if err != nil {
            ch.Close()
            conn.Close()
            logger.Error(logging.RabbitMQ, logging.ExternalService, fmt.Sprintf("Failed to declare queue '%s': %v", queueName, err), nil)
            return nil, fmt.Errorf("failed to declare queue '%s': %w", queueName, err)
        }

        service.Channels[queueName] = ch
        logger.Info(logging.RabbitMQ, logging.Startup, fmt.Sprintf("Channel created for queue: %s", queueName), nil)
    }

    return service, nil
}

// PublishMessage publishes a message to a specified queue
func (rmq *RabbitMQService) PublishMessage(queueName string, message interface{}) error {
    ch, ok := rmq.Channels[queueName]
    if !ok {
        return fmt.Errorf("queue '%s' not found", queueName)
    }

    body, err := json.Marshal(message)
    if err != nil {
        rmq.logger.Error(logging.RabbitMQ, logging.ExternalService, fmt.Sprintf("Failed to marshal message: %v", err), nil)
        return fmt.Errorf("failed to marshal message: %w", err)
    }

    err = ch.Publish(
        "",          // exchange
        queueName,   // routing key (queue name)
        false,       // mandatory
        false,       // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            Body:         body,
            DeliveryMode: amqp.Persistent, // make message persistent
        },
    )
    if err != nil {
        rmq.logger.Error(logging.RabbitMQ, logging.ExternalService, fmt.Sprintf("Failed to publish a message to queue '%s': %v", queueName, err), nil)
        return fmt.Errorf("failed to publish a message to queue '%s': %w", queueName, err)
    }

    rmq.logger.Info(logging.RabbitMQ, logging.Startup, fmt.Sprintf("Message published to queue: %s", queueName), nil)
    return nil
}

// ConsumeMessages sets up a consumer for a specified queue
func (rmq *RabbitMQService) ConsumeMessages(queueName string, handler func([]byte)) error {
    ch, ok := rmq.Channels[queueName]
    if !ok {
        return fmt.Errorf("queue '%s' not found", queueName)
    }

    msgs, err := ch.Consume(
        queueName, // queue name
        "",        // consumer name
        false,     // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if err != nil {
        rmq.logger.Error(logging.RabbitMQ, logging.ExternalService, fmt.Sprintf("Failed to register a consumer for queue '%s': %v", queueName, err), nil)
        return fmt.Errorf("failed to register a consumer for queue '%s': %w", queueName, err)
    }

    go func() {
        for d := range msgs {
            handler(d.Body) // Process the message
            d.Ack(false)    // Manually acknowledge the message after processing
        }
    }()

    rmq.logger.Info(logging.RabbitMQ, logging.Startup, fmt.Sprintf("Consuming messages from queue: %s", queueName), nil)
    return nil
}

// Close closes all channels and the RabbitMQ connection
func (rmq *RabbitMQService) Close() {
    for _, ch := range rmq.Channels {
        if err := ch.Close(); err != nil {
            rmq.logger.Error(logging.RabbitMQ, logging.ExternalService, fmt.Sprintf("Failed to close channel: %v", err), nil)
        }
    }
    if err := rmq.Connection.Close(); err != nil {
        rmq.logger.Error(logging.RabbitMQ, logging.ExternalService, fmt.Sprintf("Failed to close connection: %v", err), nil)
    }
}

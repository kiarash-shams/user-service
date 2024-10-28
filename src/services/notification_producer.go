package services

import (
	// "encoding/base64"
	"encoding/json"
	"fmt"
)

type NotificationProducer struct {
    rabbitmqService *RabbitMQService
}

func NewNotificationProducer(rabbitmqService *RabbitMQService) *NotificationProducer {
    return &NotificationProducer{
        rabbitmqService: rabbitmqService,
    }
}

func (np *NotificationProducer) QueueNotification(notification *Notification) error {
    queueName := fmt.Sprintf("%s_queue", notification.Type)
    
    notificationJSON, err := json.Marshal(notification)
   
    if err != nil {
        return fmt.Errorf("failed to marshal notification: %w", err)
    }
    // Encode to Base64
    // encodedJSON := base64.StdEncoding.EncodeToString(notificationJSON)

   
    queueName = fmt.Sprintf("%s_queue", notification.Type)

    err = np.rabbitmqService.PublishMessage(queueName, notificationJSON)
    if err != nil {
        return fmt.Errorf("failed to queue notification: %w", err)
    }

    return nil
}
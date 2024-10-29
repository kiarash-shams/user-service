package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type NotificationConsumer struct {
    rabbitmqService *RabbitMQService
    notifier        Notifier
    queueName       string
}

func NewNotificationConsumer(rabbitmqService *RabbitMQService, notificationType string) *NotificationConsumer {
    return &NotificationConsumer{
        rabbitmqService: rabbitmqService,
        notifier:        NewNotifier(notificationType),
        queueName:       fmt.Sprintf("%s_queue", notificationType),
    }
}

func (nc *NotificationConsumer) Start() error {
    return nc.rabbitmqService.ConsumeMessages(nc.queueName, nc.processNotification)
}

func (nc *NotificationConsumer) processNotification(body []byte) {
	log.Println("Received notification body:", string(body)) // چاپ محتوای JSON دریافتی
      // Remove the surrounding quotes if they exist
    bodyStr := strings.Trim(string(body), "\"")

    
    decodedBody, err := base64.StdEncoding.DecodeString(bodyStr)

    if err != nil {
        log.Printf("Failed to decode base64: %v\nReceived body: %s", err, bodyStr)
		 return
	}
   
	 log.Println("Decoded notification body:", string(decodedBody)) // چاپ محتوای رمزگشایی شده
 
    var notification Notification
    err = json.Unmarshal(decodedBody, &notification)
    if err != nil {
        log.Printf("Failed to unmarshal notification: %v", err)
        return
    }
        // Validate the notification
        if notification.Type == "" || notification.Recipient == "" || notification.Message == "" {
            log.Printf("Invalid notification: missing required fields")
            return
        }

    err = nc.notifier.Send(&notification)
    if err != nil {
        log.Printf("Failed to send notification: %v", err)
        // Here you could implement retry logic or move to a dead-letter queue
    }
}
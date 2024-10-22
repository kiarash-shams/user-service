package services

import (
	"fmt"
	"time"
	provider_sms "user-service/provider/sms"
	provider_email "user-service/provider/email"
)

// Notification struct holds the information needed for a notification
type Notification struct {
    Type      string    `json:"type"`      // sms, email, webhook
    Recipient string    `json:"recipient"`
    Message   string    `json:"message"`
    CreatedAt time.Time `json:"created_at"`
    Status    string    `json:"status"`    // pending, sent, failed
}

// Notifier interface for sending notifications
type Notifier interface {
    Send(notification *Notification) error
}

// SMSNotifier struct for sending SMS notifications
type SMSNotifier struct {
    providerManager *provider_sms.ProviderManager // Use the existing ProviderManager
}

// EmailNotifier struct for sending Email notifications
type EmailNotifier struct {
    providerManager *provider_email.EmailProviderManager
}

// NewEmailNotifier initializes the EmailNotifier
func NewEmailNotifier() *EmailNotifier {
    return &EmailNotifier{
        providerManager: provider_email.NewEmailProviderManager(),
    }
}

// NewSMSNotifier initializes the SMSNotifier
func NewSMSNotifier() *SMSNotifier {
    return &SMSNotifier{
        providerManager: provider_sms.NewProviderManager(), // Initialize the provider manager for SMS
    }
}

// Send method for SMSNotifier using ProviderManager
func (s SMSNotifier) Send(notification *Notification) error {
    err := s.providerManager.SendSMS(notification.Recipient, notification.Message)
    if err != nil {
        notification.Status = "failed"
        return err
    }
    notification.Status = "sent"
    return nil
}

// Send method for EmailNotifier
func (e EmailNotifier) Send(notification *Notification) error {
    err := e.providerManager.SendEmail(notification.Recipient, "Notification", notification.Message)
    if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }
    notification.Status = "sent"
    return nil
}

// NewNotifier returns the appropriate notifier based on the notification type
func NewNotifier(notificationType string) Notifier {
    switch notificationType {
    case "sms":
        return NewSMSNotifier() // Use the new SMS notifier
    case "email":
        return NewEmailNotifier()
    default:
        return nil
    }
}

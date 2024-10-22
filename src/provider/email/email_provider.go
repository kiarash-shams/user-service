package provider

// EmailProvider defines the interface for sending emails
type EmailProvider interface {
    SendEmail(recipient, subject, message string) error
}
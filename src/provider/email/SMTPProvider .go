package provider

import "fmt"

// SMTPProvider is a dummy implementation of the EmailProvider interface
type SMTPProvider struct{}

// SendEmail method for SMTPProvider
func (s SMTPProvider) SendEmail(recipient, subject, message string) error {
    fmt.Printf("Sending Email via SMTP to %s: %s\n", recipient, message)
    // Here you would have the logic to send an email using SMTP
    return nil // Assume success for this example
}
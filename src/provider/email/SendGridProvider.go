package provider

import "fmt"

// SendGridProvider is a dummy implementation of the EmailProvider interface
type SendGridProvider struct{}

// SendEmail method for SendGridProvider
func (s SendGridProvider) SendEmail(recipient, subject, message string) error {
    fmt.Printf("Sending Email via SendGrid to %s: %s\n", recipient, message)
    // Here you would have the logic to send an email using SendGrid
    return nil // Assume success for this example
}

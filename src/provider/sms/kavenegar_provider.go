package provider

import (
    "fmt"
)


// KavenegarProvider struct implements the SMSProvider interface
type KavenegarProvider struct{}

// SendSMS method for Kavenegar provider (dummy implementation)
func (p KavenegarProvider) SendSMS(recipient, message string) error {
    // Here you would implement the actual API call to Kavenegar
    fmt.Printf("Sending SMS via Kavenegar to %s: %s\n", recipient, message)
    return nil // return nil to indicate success
}

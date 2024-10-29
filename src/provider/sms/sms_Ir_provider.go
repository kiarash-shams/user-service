package provider

import (
    "fmt"
)


// SmsIrProvider struct implements the SMSProvider interface
type SmsIrProvider struct{}

// SendSMS method for sms.ir provider (dummy implementation)
func (p SmsIrProvider) SendSMS(recipient, message string) error {
    // Here you would implement the actual API call to sms.ir
    fmt.Printf("Sending SMS via sms.ir to %s: %s\n", recipient, message)
    return nil // return nil to indicate success
}

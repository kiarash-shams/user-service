package provider


// SMSProvider interface for SMS providers
type SMSProvider interface {
	SendSMS(recipient, message string) error
}

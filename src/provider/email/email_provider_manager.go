package provider

import (
    "fmt"
    "sync"
    "time"
)

// EmailProviderManager manages multiple Email providers
type EmailProviderManager struct {
    providers []EmailProvider
}

// NewEmailProviderManager initializes an EmailProviderManager with all available providers
func NewEmailProviderManager() *EmailProviderManager {
    return &EmailProviderManager{
        providers: []EmailProvider{
            SMTPProvider{},
            SendGridProvider{},
            // Add other providers here
        },
    }
}

// SendEmail sends the email using the first available provider
func (pm *EmailProviderManager) SendEmail(recipient, subject, message string) error {
    var wg sync.WaitGroup
    resultCh := make(chan error, len(pm.providers))
    timeout := time.After(5 * time.Second) // Set a timeout for the request

    for _, provider := range pm.providers {
        wg.Add(1)
        go func(p EmailProvider) {
            defer wg.Done()
            // Attempt to send the email and send the result (error or nil) to the result channel
            resultCh <- p.SendEmail(recipient, subject, message)
        }(provider)
    }

    // Close the result channel once all goroutines have finished
    go func() {
        wg.Wait()
        close(resultCh)
    }()

    select {
    case err := <-resultCh:
        // Return the first error received from any provider
        return err
    case <-timeout:
        return fmt.Errorf("timeout while sending email")
    }
}


package provider

import (
    "fmt"
    "sync"
    "time"
)

// ProviderManager manages multiple SMS providers
type ProviderManager struct {
    providers []SMSProvider
}

// NewProviderManager initializes a ProviderManager with all available providers
func NewProviderManager() *ProviderManager {
    return &ProviderManager{
        providers: []SMSProvider{
            SmsIrProvider{},
            KavenegarProvider{},
            // Add other providers here
        },
    }
}

// SendSMS sends the SMS using the first available provider
func (pm *ProviderManager) SendSMS(recipient, message string) error {
    var wg sync.WaitGroup
    resultCh := make(chan error, len(pm.providers))
    timeout := time.After(5 * time.Second) // Set a timeout for the request

    for _, provider := range pm.providers {
        wg.Add(1)
        go func(p SMSProvider) {
            defer wg.Done()
            resultCh <- p.SendSMS(recipient, message)
        }(provider)
    }

    go func() {
        wg.Wait()
        close(resultCh)
    }()

    select {
    case err := <-resultCh:
        return err
    case <-timeout:
        return fmt.Errorf("timeout while sending SMS")
    }
}
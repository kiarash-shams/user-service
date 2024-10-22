package provider

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// SMSIrProvider struct for integrating with SMS.ir API
type SMSIrProvider struct {
    APIKey     string
    TemplateID string
}

// NewSMSIrProvider initializes and returns an instance of SMSIrProvider
func NewSMSIrProvider(apiKey, templateID string) SMSIrProvider {
    return SMSIrProvider{
        APIKey:     apiKey,
        TemplateID: templateID,
    }
}

// SendSMS sends an SMS using the SMS.ir API
func (s SMSIrProvider) SendSMS(recipient, message string) error {
    // Prepare the request payload
    data := map[string]interface{}{
        "mobile":     recipient,
        "templateId": s.TemplateID,
        "parameters": []map[string]string{
            {"name": "PARAMETER1", "value": message}, // You may need to customize this based on the SMS.ir requirements
        },
    }

    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal request data: %w", err)
    }

    // Create the HTTP request
    req, err := http.NewRequest("POST", "https://api.sms.ir/v1/send/verify", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "text/plain")
    req.Header.Set("x-api-key", s.APIKey)

    // Create a new HTTP client with a timeout
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    // Check the response status code
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send SMS, status code: %d", resp.StatusCode)
    }

    // Optionally, you can parse the response body if needed

    fmt.Println("SMS sent successfully using SMS.ir")
    return nil
}

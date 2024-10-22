package dto

// NotificationRequest DTO for creating a notification
type NotificationRequest struct {
    Type      string `json:"type" binding:"required"`      // sms, email, webhook
    Recipient string `json:"recipient" binding:"required"`
    Message   string `json:"message" binding:"required"`
}

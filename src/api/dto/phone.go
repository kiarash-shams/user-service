package dto


type CreatePhoneRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,min=11,max=11"`
	Country string `json:"country" binding:"min=3,max=11"`
}

type UpdatePhoneRequest struct {
	MobileNumber string	`json:"mobileNumber" binding:"required,mobile,min=11,max=11"`
	Country string `json:"country" binding:"min=3,max=11"`
}

type PhoneResponse struct {
	MobileNumber string `json:"mobileNumber"`	
}
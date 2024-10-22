package dto

type CreateLocationRequest struct {
	Address   		   string `json:"address" binding:"max=255"`
	PostalCode         string `json:"postalCode" binding:"max=10"`
	City               string `json:"city" binding:"max=50"`
	Country            string `json:"country" binding:"max=50"`
	StaticPhoneNumber  string `json:"staticPhoneNumber" binding:"max=20"`
	Metadata           string `json:"metadata" binding:"max=255"`
	
}

type UpdateLocationRequest struct {
	Address   		   string `json:"address" binding:"max=255"`
	PostalCode         string `json:"postalCode" binding:"max=10"`
	City               string `json:"city" binding:"max=50"`
	Country            string `json:"country" binding:"max=50"`
	StaticPhoneNumber  string `json:"staticPhoneNumber" binding:"max=20"`
	Metadata           string `json:"metadata" binding:"max=255"`
}

type LocationResponse struct {
	Address   		   string `json:"address"`
	PostalCode         string `json:"postalCode"`
	City               string `json:"city"`
	Country            string `json:"country"`
	StaticPhoneNumber  string `json:"staticPhoneNumber"`
	Metadata           string `json:"metadata"`
	Verification       string `json:"verification"`
}

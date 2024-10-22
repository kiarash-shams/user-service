package dto

type CreateLabelRequest struct {
	Key 		   string `json:"key" binding:"required,min=3,max=255"`	
	Value 		   string `json:"value" binding:"required,min=3,max=255"`	
	Scope 		   string `json:"scope" binding:"min=3,max=255"`	
	Description   string `json:"description" binding:"required,max=255"`
}

type UpdateLabelRequest struct {
	Key 		   string `json:"key" binding:"min=3,max=255"`	
	Value 		   string `json:"value" binding:"min=3,max=255"`	
	Description   string `json:"description" binding:"max=255"`
}

type LabelResponse struct {
	Key 		   string `json:"key"`	
	Value 		   string `json:"value"`	
	Scope 		   string `json:"scope"`	
	Description   string `json:"description"`
}
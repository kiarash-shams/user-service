package dto

type CreateProfileRequest struct {
	FirstNameEncrypted string `json:"firstName" binding:"required,min=3,max=255"`	
	LastNameEncrypted  string `json:"lastName" binding:"required,min=3,max=255"`	
	FatherName 		   string `json:"fatherName" binding:"min=3,max=255"`	
	NidEncrypted 	   string `json:"nid" binding:"required,min=3,max=11"`
	DobEncrypted 	   string `json:"dob" binding:"required,max=255"`	
	Metadata           string `json:"metadata" binding:"max=255"`
}

type UpdateProfileRequest struct {
	FirstNameEncrypted string `json:"firstName" binding:"min=3,max=255"`	
	LastNameEncrypted  string `json:"lastName" binding:"min=3,max=255"`	
	FatherName 		   string `json:"fatherName" binding:"min=3,max=255"`	
	NidEncrypted 	   string `json:"nid" binding:"min=3,max=11"`
	DobEncrypted 	   string `json:"dob" binding:"max=255"`	
	Metadata           string `json:"metadata" binding:"max=255"`
}

type ProfileResponse struct {
	FirstNameEncrypted 	string `json:"firstName"`	
	LastNameEncrypted 	string `json:"lastName"`	
	FatherName 		    string `json:"fatherName"`	
	NidEncrypted 		string `json:"nid"`
	DobEncrypted 		string `json:"dob"`	
	Metadata           	string `json:"metadata"`
	Verification       	string `json:"verification"`
}

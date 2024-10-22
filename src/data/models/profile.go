package model

import (
	"database/sql"
	"time"
)


type Profile struct {
	Id 				   	int 	`gorm:"primarykey"`
	UserId             	int    	`gorm:"not null;unique"`
	NidEncrypted        string 	`gorm:"type:string;size:255;null;unique"`    // National ID
	FirstNameEncrypted  string 	`gorm:"type:string;size:255;null"`   //  first name
	LastNameEncrypted   string 	`gorm:"type:string;size:255;null"`   //  last name
	FatherName  	   	string 	`gorm:"type:string;size:255;null"` 
	DobEncrypted       	string 	`gorm:"type:string;size:255;null"`   //  date of birth
	Metadata           	string 	`gorm:"type:text;null;default:null"` // Additional metadata as text
	Verification       	bool   	`gorm:"default:false"`               // Verification status	
	
	CreatedAt  		   	time.Time    `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt 		   	sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	ModifiedBy 		   	*sql.NullInt64 `gorm:"null"`
	DeletedAt  		   	sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedBy  		   	*sql.NullInt64 `gorm:"null"`

	User               	User   `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}
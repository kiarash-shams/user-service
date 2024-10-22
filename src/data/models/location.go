package model

import (
	"database/sql"
	"time"
)


type Location struct {
	Id         			int    `gorm:"primarykey"`
	UserId     			int    `gorm:"not null;unique"`
	PostalCode 			string `gorm:"type:string;size:20;null;unique"` // Postal code (unique)
	Address    			string `gorm:"type:string;size:255;null"`      // Address
	City       			string `gorm:"type:string;size:50;null"`       // City
	Country    			string `gorm:"type:string;size:50;null"`       // Country
	StaticPhoneNumber 	string `gorm:"type:string;size:20;null;unique"`  // Phone number (unique)
	Metadata   			string `gorm:"type:text;null;default:null"`    // Additional metadata as text
	Verification       	bool   `gorm:"default:false"`               // Verification status	
	
	CreatedAt   		time.Time     `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt  		sql.NullTime  `gorm:"type:TIMESTAMP with time zone;null"`
	ModifiedBy  		*sql.NullInt64 `gorm:"null"`
	DeletedAt   		sql.NullTime  `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedBy   		*sql.NullInt64 `gorm:"null"`

	User User `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}
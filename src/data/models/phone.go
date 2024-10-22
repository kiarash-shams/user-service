package model

import (
	"database/sql"
	"time"
)

type Phone struct {
	Id 		    	int 			`gorm:"primarykey"` 
	UserId        	int 			`gorm:"not null"`
	Country       	string 			`gorm:"type:character;size:10;null"`
	Number			string  		`gorm:"type:character varying;not null;unique"`
	Verification    bool   			`gorm:"default:false"`
	ValidatedAt   	sql.NullTime 	`gorm:"type:TIMESTAMP with time zone;null"`
	
	CreatedAt  		time.Time    	`gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt 		sql.NullTime 	`gorm:"type:TIMESTAMP with time zone;null"`
	ModifiedBy 		*sql.NullInt64 	`gorm:"null"`
	DeletedAt  		sql.NullTime 	`gorm:"type:TIMESTAMP with time zone;null"`
	DeletedBy  		*sql.NullInt64 	`gorm:"null"`

	User          	User  `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"` // Foreign key to User table
}
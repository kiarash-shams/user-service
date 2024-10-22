package model



import (
	"database/sql"
	"time"
)


type Document struct {
	Id				int 	`gorm:"primarykey"`
	UserId			int     `gorm:"not null"`
	DocCategory		string  `gorm:"type:character varying;null"`

	Name			string `gorm:"size:100;type:string;not null"`
	Directory		string `gorm:"size:100;type:string;not null"`
	Description		string `gorm:"size:500;type:string;not null"`
	MimeType		string `gorm:"size:20;type:string;not null"`
	
	CreatedAt		time.Time    	`gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt		sql.NullTime	`gorm:"type:TIMESTAMP with time zone;null"`
	ModifiedBy		*sql.NullInt64 	`gorm:"null"`
	DeletedAt		sql.NullTime 	`gorm:"type:TIMESTAMP with time zone;null"`
	DeletedBy		*sql.NullInt64 	`gorm:"null"`


	User            User      `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"` // Foreign key to User table
	
}

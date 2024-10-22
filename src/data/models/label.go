package model

import (
	"database/sql"
	"time"
)

type Label struct {
	Id 		    int    `gorm:"primarykey"`
	UserId      int    `gorm:"not null"`
	Key         string `gorm:"type:character varying;size:255;not null"`
	Value       string `gorm:"type:character varying;size:255;not null"`
	Scope       string `gorm:"type:character varying;size:255;default:'public'"`
	Description string `gorm:"type:character varying;size:255;null"`
	
	CreatedAt  		   time.Time    `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt 		   sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	ModifiedBy 		   *sql.NullInt64 `gorm:"null"`
	DeletedAt  		   sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedBy  		   *sql.NullInt64 `gorm:"null"`

	User               User   `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

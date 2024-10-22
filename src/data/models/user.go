package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id           int    `gorm:"primarykey"`
	UID          string `gorm:"type:character varying;size:20;not null;unique"` // Unique identifier for the user
	MobileNumber string `gorm:"type:string;size:11;null;unique;default:null"`   // User's mobile number (nullable, unique)
	Email        string `gorm:"type:string;size:64;null;unique;default:null"`   // User's email address (nullable, unique)
	Password     string `gorm:"type:string;size:64;not null"`                   // User's password (hashed)
	Level        int    `gorm:"default:0"`                                      // User's level, default is 0
	Otp          bool   `gorm:"null;default:false"`                             // OTP status (nullable, default is false)
	State        bool   `gorm:"default:true"`                                   // Account state (true for active, false for inactive)
	Username     string `gorm:"type:string;size:20;null;unique"`                // Username (unique)
	ReferralId   int    `gorm:"null"`                                           // Referral ID (nullable)

	CreatedAt  time.Time      `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime   `gorm:"type:TIMESTAMP with time zone;null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedAt  sql.NullTime   `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`

	UserRoles *[]UserRole // Roles associated with the user
	Phones    []Phone     // Associated phone numbers (relationship to Phone model)
	Profile   []Profile   // Associated profiles for the user (relationship to Profile model)
	Label     []Label     // Associated labels for the user (relationship to Label model)
	Document  []Document  // Associated documents for the user (relationship to Document model)
}

type Role struct {
	BaseModel
	Name      string `gorm:"type:string;size:10;not null,unique"`
	UserRoles *[]UserRole
}

type UserRole struct {
	BaseModel
	User   User `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	Role   Role `gorm:"foreignKey:RoleId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	UserId int
	RoleId int
}

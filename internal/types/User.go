package types

import (
	"time"
)

type Example struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `gorm:"type:varchar(255)" json:"first_name"`
	LastName  string `gorm:"type:varchar(255)" json:"last_name"`
}

// IPAddress represents the ip_addresses table.
type UserIPAddress struct {
	ID         uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	IP1        string          `gorm:"type:varchar(15);not null" json:"ip1"`
	IP2        string          `gorm:"type:varchar(15)" json:"ip2"`
	IP3        string          `gorm:"type:varchar(15)" json:"ip3"`
	IP4        string          `gorm:"type:varchar(15)" json:"ip4"`
	LastUsed   *time.Time      `gorm:"type:timestamp" json:"last_used"`
	Users      []User          `gorm:"foreignKey:IPAddressID" json:"users,omitempty"`
	SiteVisits []UserSiteVisit `gorm:"foreignKey:IPAddressID" json:"site_visits,omitempty"`
}

// User represents the users table.
type User struct {
	ID            uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName      string                `gorm:"type:varchar(255)" json:"full_name"`
	Nickname      string                `gorm:"type:varchar(50)" json:"nickname"`
	Email         string                `gorm:"type:varchar(255);not null;unique" json:"email"`
	EmailVerified bool                  `gorm:"default:false" json:"email_verified"`
	PasswordID    uint                  `gorm:"not null" json:"password_id"`
	CreatedAt     time.Time             `gorm:"autoCreateTime" json:"created_at"`
	IPAddressID   *uint                 `json:"ip_address_id,omitempty"`
	Addresses     []UserPhysicalAddress `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	SiteVisits    []UserSiteVisit       `gorm:"foreignKey:UserID" json:"site_visits,omitempty"`
}

// Password represents the passwords table.
type UserPassword struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint       `gorm:"not null;unique" json:"user_id"`
	HashedPassword  string     `gorm:"type:varchar(255);not null" json:"hashed_password"`
	ResetInProgress bool       `gorm:"default:false" json:"reset_in_progress"`
	ResetCode       string     `gorm:"type:varchar(255)" json:"reset_code"`
	ResetExpires    *time.Time `gorm:"type:timestamp" json:"reset_expires"`
	Active          bool       `gorm:"default:true" json:"active"`
	User            User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
}

// PhysicalAddress represents the physical_addresses table.
type UserPhysicalAddress struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint   `gorm:"not null" json:"user_id"`
	StreetNumber int    `gorm:"not null" json:"street_number"`
	Directional  string `gorm:"type:varchar(10)" json:"directional"`
	Street       string `gorm:"type:varchar(255);not null" json:"street"`
	Suffix       string `gorm:"type:varchar(50)" json:"suffix"`
	UnitType     string `gorm:"type:varchar(50)" json:"unit_type"`
	UnitNumber   int    `json:"unit_number"`
	ZipCode      string `gorm:"type:varchar(20);not null" json:"zip_code"`
	CountryCode  string `gorm:"type:varchar(5);not null" json:"country_code"`
	Primary      bool   `gorm:"default:false" json:"primary"`
	Active       bool   `gorm:"default:true" json:"active"`
	IsBilling    bool   `gorm:"default:false" json:"is_billing"`
	IsShipping   bool   `gorm:"default:false" json:"is_shipping"`
	User         User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
}

// SiteVisit represents the site_visits table.
type UserSiteVisit struct {
	ID                   uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID               uint          `gorm:"not null" json:"user_id"`
	IPAddressID          *uint         `json:"ip_address_id,omitempty"`
	VisitStart           time.Time     `gorm:"type:timestamp;not null" json:"visit_start"`
	VisitLastInteraction *time.Time    `gorm:"type:timestamp" json:"visit_last_interaction"`
	ReferrerURL          string        `gorm:"type:varchar(255)" json:"referrer_url"`
	User                 User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
	IPAddress            UserIPAddress `gorm:"foreignKey:IPAddressID" json:"ip_address"`
}

// EmailVerification represents the email_verifications table.
type UserEmailVerification struct {
	Token     string    `gorm:"primaryKey;type:varchar(255)" json:"token"`
	Email     string    `gorm:"type:varchar(255);not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

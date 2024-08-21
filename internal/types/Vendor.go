package types

import (
	"time"

	"gorm.io/gorm"
)

type VendorIPAddress struct {
	ID       uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	IP1      string     `gorm:"type:varchar(15);not null" json:"ip1"`
	IP2      string     `gorm:"type:varchar(15)" json:"ip2"`
	IP3      string     `gorm:"type:varchar(15)" json:"ip3"`
	IP4      string     `gorm:"type:varchar(15)" json:"ip4"`
	LastUsed *time.Time `gorm:"type:timestamp" json:"last_used"`

	// Associations
	Vendors    []Vendor          `gorm:"foreignKey:IPAddressID" json:"vendors,omitempty"`
	SiteVisits []VendorSiteVisit `gorm:"foreignKey:IPAddressID" json:"site_visits,omitempty"`
}

type Vendor struct {
	gorm.Model
	TenantID    string `gorm:"type:varchar(255);not null" json:"tenant_id"`
	Logo        string `gorm:"type:varchar(255)" json:"logo"`
	CompanyName string `gorm:"type:varchar(255);not null" json:"company_name"`
	TradingName string `gorm:"type:varchar(255);not null" json:"trading_name"`

	PhoneNo       string `gorm:"type:varchar(20)" json:"phone_no"`
	Email         string `gorm:"type:varchar(255);not null;unique" json:"email"`
	EmailVerified bool   `gorm:"default:false" json:"email_verified"`
	PasswordID    uint   `gorm:"not null" json:"password_id"`
	IPAddressID   *uint  `json:"ip_address_id,omitempty"`
	TryCount      int    `gorm:"default:0" json:"try_count"`

	// Associations
	OTPs       []VendorOTP             `gorm:"foreignKey:VendorID" json:"otps,omitempty"`
	Passwords  []VendorPassword        `gorm:"foreignKey:VendorID" json:"passwords,omitempty"`
	Addresses  []VendorPhysicalAddress `gorm:"foreignKey:VendorID" json:"addresses,omitempty"`
	SiteVisits []VendorSiteVisit       `gorm:"foreignKey:VendorID" json:"site_visits,omitempty"`
}

type VendorPassword struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	VendorID        uint       `gorm:"not null;unique" json:"vendor_id"`
	HashedPassword  string     `gorm:"type:varchar(255);not null" json:"hashed_password"`
	ResetInProgress bool       `gorm:"default:false" json:"reset_in_progress"`
	ResetCode       string     `gorm:"type:varchar(255)" json:"reset_code"`
	ResetExpires    *time.Time `gorm:"type:timestamp" json:"reset_expires"`
	Active          bool       `gorm:"default:true" json:"active"`
	Vendor          Vendor     `gorm:"foreignKey:VendorID;constraint:OnDelete:CASCADE;" json:"vendor"`
}
type VendorOTP struct {
	gorm.Model
	VendorID  uint   `gorm:"not null" json:"vendor_id"`
	OTP       string `gorm:"type:varchar(6);not null" json:"otp"`
	Revoked   bool   `gorm:"default:false" json:"revoked"`
	ExpiresAt time.Time

	// Associations
	Vendor Vendor `gorm:"foreignKey:VendorID" json:"vendor"`
}

type VendorPhysicalAddress struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	VendorID     uint   `gorm:"not null" json:"vendor_id"`
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
}

type VendorSiteVisit struct {
	ID          uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	VendorID    uint            `gorm:"not null" json:"vendor_id"`
	IPAddressID uint            `gorm:"not null" json:"ip_address_id"`
	IPAddress   VendorIPAddress `gorm:"foreignKey:IPAddressID" json:"ip_address"`
	VisitTime   time.Time       `gorm:"autoCreateTime" json:"visit_time"`
}

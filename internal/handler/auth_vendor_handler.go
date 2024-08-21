package handler

import (
	"net/http"
	"time"

	"github.com/Satishcg12/multicommers/internal/types"
	"github.com/Satishcg12/multicommers/utils/password"
	randomString "github.com/Satishcg12/multicommers/utils/string"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type (
	AuthVendorHandler struct {
	}
	AuthVendorHandlerInterface interface {
		Register(c echo.Context) error
		VerifyOTP(c echo.Context) error
		ResendOTP(c echo.Context) error
		Login(c echo.Context) error
		RequestResetPassword(c echo.Context) error
		ResetPassword(c echo.Context) error
		Logout(c echo.Context) error
	}
	registerRequest struct {
		CompanyName string `json:"company_name" form:"company_name" query:"company_name" validate:"required,min=3,max=255"`
		TradingName string `json:"trading_name" form:"trading_name" query:"trading_name" validate:"required,min=3,max=255"`
		Email       string `json:"email" form:"email" query:"email" validate:"required,email"`
		Password    string `json:"password" form:"password" query:"password" validate:"required,password,min=8,max=50"`
		ConfirmPass string `json:"confirm_password" form:"confirm_password" query:"confirm_password" validate:"required,eqfield=Password"`
		PhoneNumber string `json:"phone_number" form:"phone_number" query:"phone_number" validate:"required,min=10,max=15"`
	}
	verifyOTPRequest struct {
		Email string `json:"email" form:"email" query:"email" validate:"required,email"`
		OTP   string `json:"otp" form:"otp" query:"otp" validate:"required"`
	}
	ResendOTPRequest struct {
		Email string `json:"email" form:"email" query:"email" validate:"required,email"`
	}

	loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	RequestResetPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}
	ResetPasswordRequest struct {
		Email           string `json:"email" validate:"required,email"`
		Token           string `json:"token" validate:"required"`
		Password        string `json:"password" validate:"required,password,min=8,max=50"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	}
)

func NewAuthVendorHandler() AuthVendorHandlerInterface {
	return &AuthVendorHandler{}
}

func (h *AuthVendorHandler) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.(*echo.HTTPError).Message)
	}

	db := c.Get("db").(*gorm.DB)

	// check if email already exists
	if err := db.Where("email = ?", req.Email).First(&types.Vendor{}).Error; err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email already exists"})
	}
	// check if company name already exists
	if err := db.Where("company_name = ?", req.CompanyName).First(&types.Vendor{}).Error; err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "company name already exists"})
	}
	// check if trading name already exists
	if err := db.Where("trading_name = ?", req.TradingName).First(&types.Vendor{}).Error; err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "trading name already exists"})
	}

	// hash password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error hashing password"})
	}

	// save to db with transaction
	tx := db.Begin()
	// create vendor
	vendor := types.Vendor{
		CompanyName: req.CompanyName,
		TradingName: req.TradingName,
		Email:       req.Email,
		PhoneNo:     req.PhoneNumber,
	}
	if err := tx.Create(&vendor).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error creating vendor"})
	}
	// create password
	password := types.VendorPassword{
		VendorID:       vendor.ID,
		HashedPassword: hashedPassword,
	}
	if err := tx.Create(&password).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error creating password"})
	}
	otp := types.VendorOTP{
		VendorID:  vendor.ID,
		OTP:       randomString.GenerateRandomString(6),
		ExpiresAt: time.Now().Add(time.Minute * 15),
		Revoked:   false,
	}

	if err := tx.Create(&otp).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error creating otp"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error creating vendor"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "success"})

}

func (h *AuthVendorHandler) VerifyOTP(c echo.Context) error {
	var req verifyOTPRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.(*echo.HTTPError).Message)
	}

	db := c.Get("db").(*gorm.DB)

	// check if email exists
	vendor := types.Vendor{}
	if err := db.Where("email = ?", req.Email).First(&vendor).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email does not exist"})
	}

	// check if otp exists
	otp := types.VendorOTP{}
	if err := db.Where("vendor_id = ? AND otp = ? AND revoked = false AND expires_at > ?", vendor.ID, req.OTP, time.Now()).First(&otp).Error; err != nil {
		vendor.TryCount++
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid otp"})
	}
	if vendor.TryCount >= 3 {
		otp.Revoked = true
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "too many attempts"})
	}

	// update vendor
	if err := db.Model(&vendor).Update("email_verified", true).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error updating vendor"})
	}

	// update otp
	return c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func (h *AuthVendorHandler) ResendOTP(c echo.Context) error {
	var req ResendOTPRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.(*echo.HTTPError).Message)
	}

	db := c.Get("db").(*gorm.DB)

	// check if email exists
	vendor := types.Vendor{}
	if err := db.Where("email = ?", req.Email).First(&vendor).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email does not exist"})
	}

	// check if otp exists
	otp := types.VendorOTP{}
	if err := db.Where("vendor_id = ? AND revoked = false AND expires_at > ?", vendor.ID, time.Now()).First(&otp).Error; err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "otp already sent"})
	}
	if otp.CreatedAt.Add(time.Minute * 1).After(time.Now()) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "otp already sent"})
	}
	otp.Revoked = true
	if err := db.Model(&otp).Update("revoked", true).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error updating otp"})
	}

	// create otp
	otp = types.VendorOTP{
		VendorID:  vendor.ID,
		OTP:       randomString.GenerateRandomString(6),
		ExpiresAt: time.Now().Add(time.Minute * 15),
		Revoked:   false,
	}

	if err := db.Create(&otp).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error creating otp"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func (h *AuthVendorHandler) Login(c echo.Context) error {
	return nil
}

func (h *AuthVendorHandler) RequestResetPassword(c echo.Context) error {
	return nil
}

func (h *AuthVendorHandler) ResetPassword(c echo.Context) error {
	return nil
}

func (h *AuthVendorHandler) Logout(c echo.Context) error {
	return nil
}

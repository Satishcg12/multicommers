package handler

import (
	"github.com/labstack/echo/v4"
)

type (
	AuthHandler struct {
	}
	AuthHandlerInterface interface {
		Register(c echo.Context) error
		VerifyEmail(c echo.Context) error
		Login(c echo.Context) error
		RequestResetPassword(c echo.Context) error
		ResetPassword(c echo.Context) error
		Logout(c echo.Context) error
	}
	registerRequest struct {
		Email    string `json:"email" form:"email" query:"email" validate:"required,email"`
		Password string `json:"password" form:"password" query:"password" validate:"required,password,min=8,max=50"`
		FullName string `json:"full_name" form:"full_name" query:"full_name" validate:"required"`
	}
	verifyEmailRequest struct {
		Email string `json:"email" form:"email" query:"email" validate:"required,email"`
		Token string `json:"token" form:"token" query:"token" validate:"required"`
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

func NewAuthHandler() AuthHandlerInterface {
	return &AuthHandler{}
}

func (h *AuthHandler) Register(c echo.Context) error {
	return nil
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	return nil
}

func (h *AuthHandler) Login(c echo.Context) error {
	return nil
}

func (h *AuthHandler) RequestResetPassword(c echo.Context) error {
	return nil
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	return nil
}

func (h *AuthHandler) Logout(c echo.Context) error {
	return nil
}


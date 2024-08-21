package routes

import (
	"github.com/Satishcg12/multicommers/internal/handler"
	"github.com/labstack/echo/v4"
)

// RegisterVendorAuthRoutes function
func RegisterVendorAuthRoutes(e *echo.Group) {
	h := handler.NewAuthVendorHandler()

	g := e.Group("/auth/vendor")
	{
		g.POST("/register", h.Register)
		g.POST("/verify-otp", h.VerifyOTP)
		g.POST("/login", h.Login)
		g.POST("/request-reset-password", h.RequestResetPassword)
		g.POST("/reset-password", h.ResetPassword)
		g.POST("/logout", h.Logout)
	}

}

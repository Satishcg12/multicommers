package routes

import (
	"github.com/Satishcg12/multicommers/internal/handler"
	"github.com/labstack/echo/v4"
)

// RegisterAuthRoutes function
func RegisterUserAuthRoutes(e *echo.Group) {
	h := handler.NewAuthHandler()

	g := e.Group("/auth/user")
	{
		g.POST("/register", h.Register)
		g.POST("/verify-email", h.VerifyEmail)
		g.POST("/login", h.Login)
		g.POST("/request-reset-password", h.RequestResetPassword)
		g.POST("/reset-password", h.ResetPassword)
		g.POST("/logout", h.Logout)
	}

}

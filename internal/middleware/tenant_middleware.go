package middleware

import (
	"net/http"
	"strings"

	"github.com/Satishcg12/multicommers/internal/database"
	"github.com/labstack/echo/v4"
)

func TenantDBMiddleware(dbManager *database.DatabaseManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the tenant ID (e.g., from subdomain or header)
			host := c.Request().Host
			parts := strings.Split(host, ".")
			if len(parts) < 3 {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid subdomain")
			}
			tenantID := parts[0]

			// Get the corresponding DB connection
			dbConn, err := dbManager.GetDB(tenantID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to tenant database")
			}

			// Set the database connection in the context
			c.Set("db", dbConn)

			return next(c)
		}
	}
}

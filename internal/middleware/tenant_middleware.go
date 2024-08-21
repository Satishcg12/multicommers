package middleware

import (
	"net/http"
	"strings"

	"github.com/Satishcg12/multicommers/internal/database"
	"github.com/Satishcg12/multicommers/utils/dotenv"
	"github.com/labstack/echo/v4"
)

func TenantDBMiddleware(dbManager *database.DatabaseManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the tenant ID (e.g., from subdomain or header)
			subdomain := c.Request().Header.Get("X-Tenant-ID")

			if subdomain == "" {
				subdomain = strings.Split(c.Request().Host, ".")[0]
			}
			if subdomain == "" {
				subdomain = dotenv.GetEnvOrDefault("DB_NAME", "multicommers")
			}

			// Get the tenant database
			db, err := dbManager.GetDB(subdomain)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}

			// Set the database connection on the context
			c.Set("db", db)

			// Continue processing the request
			return next(c)

		}
	}
}

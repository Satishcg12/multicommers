package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/Satishcg12/multicommers/internal/database"
	myMiddleware "github.com/Satishcg12/multicommers/internal/middleware"
	"github.com/Satishcg12/multicommers/internal/router"
	"github.com/Satishcg12/multicommers/internal/types"
	"github.com/Satishcg12/multicommers/utils/dotenv"
	"github.com/Satishcg12/multicommers/utils/email"
	"github.com/Satishcg12/multicommers/utils/validators"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	TenantManager *database.DatabaseManager
	MailServer    *email.EmailDaemon
)

type (
	ServerConfig struct {
		Host string
		Port string
	}
	Server struct {
		config ServerConfig
		e      *echo.Echo
	}

	ServerInterface interface {
		Start() error
	}
)

func NewServer(config ServerConfig) ServerInterface {
	e := echo.New()
	return &Server{
		config: config,
		e:      e,
	}
}

func (s *Server) Start() error {

	// tenant manager
	tenantManager := database.NewDatabaseManager(5 * time.Minute)

	// set up main db
	err := tenantManager.InitMainDB(
		types.VendorIPAddress{},
		types.Vendor{},
		types.VendorPassword{},
		types.VendorPhysicalAddress{},
		types.VendorSiteVisit{},
	)
	if err != nil {
		log.Fatalf("Error initializing main db: %s", err)
	}

	// set tenant manager
	TenantManager = tenantManager

	// connect to mail server
	mailServer := email.NewEmailDaemon(
		dotenv.GetEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		dotenv.GetEnvOrDefault("SMTP_PORT", "587"),
		dotenv.GetEnvOrDefault("SMTP_USERNAME", ""),
		dotenv.GetEnvOrDefault("SMTP_PASSWORD", ""),
	)
	mailServer.Start()

	// middlewares
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
	s.e.Use(myMiddleware.TenantDBMiddleware(tenantManager))

	// custom validator
	s.e.Validator = validators.NewValidator()

	// init routes
	router.Init(s.e)

	// init server
	address := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	log.Printf("Server is running at %s", address)

	return s.e.Start(address)

}

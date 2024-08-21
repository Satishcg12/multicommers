package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Satishcg12/multicommers/utils/dotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type tenantInfo struct {
	db         *gorm.DB
	lastAccess time.Time
	closeChan  chan struct{}
}

type DatabaseManager struct {
	mu      sync.Mutex
	tenants map[string]*tenantInfo
	timeout time.Duration
}

func NewDatabaseManager(timeout time.Duration) *DatabaseManager {
	return &DatabaseManager{
		tenants: make(map[string]*tenantInfo),
		timeout: timeout,
	}
}

func (manager *DatabaseManager) GetDB(tenantID string) (*gorm.DB, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	// If the connection already exists, return it and update the last access time
	if tenant, exists := manager.tenants[tenantID]; exists {
		// Reset the timeout
		tenant.closeChan <- struct{}{}
		return tenant.db, nil
	}

	// Create a new database connection for the tenant
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		dotenv.GetEnvOrDefault("DB_HOST", "localhost"),
		dotenv.GetEnvOrDefault("DB_PORT", "5432"),
		dotenv.GetEnvOrDefault("DB_USERNAME", "root"),
		dotenv.GetEnvOrDefault("DB_PASSWORD", ""),
		tenantID,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create a close channel and start the timeout goroutine
	closeChan := make(chan struct{})
	tenant := &tenantInfo{
		db:         db,
		lastAccess: time.Now(),
		closeChan:  closeChan,
	}
	manager.tenants[tenantID] = tenant
	go manager.startTimeout(tenantID, tenant)

	return db, nil
}

func (manager *DatabaseManager) startTimeout(tenantID string, tenant *tenantInfo) {
	for {
		select {
		case <-time.After(manager.timeout):
			manager.mu.Lock()
			if time.Since(tenant.lastAccess) > manager.timeout {
				sqlDB, err := tenant.db.DB()
				if err == nil {
					sqlDB.Close()
				}
				delete(manager.tenants, tenantID)
			}
			manager.mu.Unlock()
			return
		case <-tenant.closeChan:
			// Reset the timeout
			manager.mu.Lock()
			tenant.lastAccess = time.Now()
			manager.mu.Unlock()
		}
	}
}

func (manager *DatabaseManager) AddTenant(tenantID string, models ...interface{}) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	// Create the tenant's database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		dotenv.GetEnvOrDefault("DB_HOST", "localhost"),
		dotenv.GetEnvOrDefault("DB_PORT", "5432"),
		dotenv.GetEnvOrDefault("DB_USERNAME", "root"),
		dotenv.GetEnvOrDefault("DB_PASSWORD", ""),
		"postgres",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Check if the tenant already exists
	if err := db.Exec("SELECT 1 FROM pg_database WHERE datname = ?", tenantID).Error; err == nil {
		return fmt.Errorf("tenant %s already exists", tenantID)
	}
	// Create the tenant's database
	createDB := fmt.Sprintf("CREATE DATABASE %s", tenantID)
	if err := db.Exec(createDB).Error; err != nil {
		return err
	}

	// Connect to the new tenant's database
	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		dotenv.GetEnvOrDefault("DB_HOST", "localhost"),
		dotenv.GetEnvOrDefault("DB_PORT", "5432"),
		dotenv.GetEnvOrDefault("DB_USERNAME", "root"),
		dotenv.GetEnvOrDefault("DB_PASSWORD", ""),
		tenantID,
	)

	tenantDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// AutoMigrate the tables
	if err := tenantDB.AutoMigrate(models...); err != nil {
		return err
	}

	// Store the connection and last access time in the map
	closeChan := make(chan struct{})
	tenant := &tenantInfo{
		db:         tenantDB,
		lastAccess: time.Now(),
		closeChan:  closeChan,
	}
	manager.tenants[tenantID] = tenant
	go manager.startTimeout(tenantID, tenant)

	return nil
}

func (manager *DatabaseManager) DeleteTenant(tenantID string) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	// Close the tenant's database connection
	if tenant, exists := manager.tenants[tenantID]; exists {
		sqlDB, err := tenant.db.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
		delete(manager.tenants, tenantID)
	}

	// Drop the tenant's database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		dotenv.GetEnvOrDefault("DB_HOST", "localhost"),
		dotenv.GetEnvOrDefault("DB_PORT", "5432"),
		dotenv.GetEnvOrDefault("DB_USERNAME", "root"),
		dotenv.GetEnvOrDefault("DB_PASSWORD", ""),
		"postgres",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// check if the tenant exists
	if err := db.Exec("SELECT 1 FROM pg_database WHERE datname = ?", tenantID).Error; err != nil {
		return fmt.Errorf("tenant %s does not exist", tenantID)
	}

	// Drop the tenant's database
	dropDB := fmt.Sprintf("DROP DATABASE %s", tenantID)
	if err := db.Exec(dropDB).Error; err != nil {
		return err
	}

	return nil
}

// set up main db

func (manager *DatabaseManager) InitMainDB(models ...interface{}) error {
	// Create the main database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		dotenv.GetEnvOrDefault("DB_HOST", "localhost"),
		dotenv.GetEnvOrDefault("DB_PORT", "5432"),
		dotenv.GetEnvOrDefault("DB_USERNAME", "root"),
		dotenv.GetEnvOrDefault("DB_PASSWORD", ""),
		"postgres",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// check if the main database already exists
	if err := db.Exec("SELECT 1 FROM pg_database WHERE datname = ?", dotenv.GetEnvOrDefault("DB_NAME", "multicommers")).Error; err != nil {

		// Create the main database
		createDB := fmt.Sprintf("CREATE DATABASE %s", dotenv.GetEnvOrDefault("DB_NAME", "multicommers"))
		if err := db.Exec(createDB).Error; err != nil {
			return err
		}
	}

	// Connect to the main database
	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		dotenv.GetEnvOrDefault("DB_HOST", "localhost"),
		dotenv.GetEnvOrDefault("DB_PORT", "5432"),
		dotenv.GetEnvOrDefault("DB_USERNAME", "root"),
		dotenv.GetEnvOrDefault("DB_PASSWORD", ""),
		dotenv.GetEnvOrDefault("DB_NAME", "multicommers"),
	)

	mainDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// AutoMigrate the tables
	if err := mainDB.AutoMigrate(models...); err != nil {
		return err
	}
	log.Println("Main database initialized")

	return nil
}

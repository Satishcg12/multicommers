package database

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	ConfigDatabase struct {
		Host     string
		Port     string
		Username string
		Password string
		Database string
	}
	Database struct {
		config ConfigDatabase
	}
	DatabaseInterface interface {
		Connect() (*gorm.DB, error)
		AutoMigrate(db *gorm.DB) error
	}
)

func NewDatabase(config ConfigDatabase) DatabaseInterface {
	return &Database{
		config: config,
	}
}

func (d *Database) Connect() (*gorm.DB, error) {
	// postgres connection string
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", d.config.Host, d.config.Port, d.config.Username, d.config.Password, d.config.Database)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	// if flag -migration is set to true
	var migration bool
	flag.BoolVar(&migration, "m", false, "Auto migrate database")
	flag.Parse()
	if migration {
		err = d.AutoMigrate(db)
		if err != nil {
			return nil, err
		}
		log.Println("Database migrated")
	}
	return db, nil
}

func (d *Database) AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
	// &types.User{},
	// &types.Category{},
	// &types.Product{},
	// &types.Images{},
	// &types.ProductAttribute{},
	// &types.ProductSKU{},
	// &types.Inventory{},
	// &types.Wishlist{},
	// &types.Review{},
	// &types.Cart{},
	// &types.CartItem{},
	// &types.Address{},
	// &types.OrderDetail{},
	// &types.OrderItem{},
	// &types.PaymentDetail{},
	)
	if err != nil {
		return err
	}
	return nil

}

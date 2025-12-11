package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"inventory-api/internal/config"
	"inventory-api/internal/models"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) {
	dsn := ""
	
	if cfg.DBPassword == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s dbname=%s port=%s sslmode=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
		)
	}
	
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	fmt.Println("Database connected successfully!")
	
	MigrateDB()
}

func MigrateDB() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.ActivityLog{},
	)
	
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	fmt.Println("Database migration completed!")
}
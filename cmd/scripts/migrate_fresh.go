package main

import (
	"log"

	"github.com/joho/godotenv"

	"inventory-api/internal/config"
	"inventory-api/internal/database"
	"inventory-api/internal/models"

	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
	
	cfg := config.LoadConfig()
	
	database.ConnectDB(cfg)
	
	log.Println("Starting fresh migration...")
	
	err := database.DB.Migrator().DropTable(
		&models.User{},
		&models.Item{},
		&models.ActivityLog{},
	)
	if err != nil {
		log.Fatal("Failed to drop tables:", err)
	}
	log.Println("Tables dropped successfully!")
	
	err = database.DB.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.ActivityLog{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	log.Println("Fresh migration completed successfully!")
	
	seedSampleData(database.DB)
}

func seedSampleData(db *gorm.DB) {
	log.Println("Seeding sample data...")
	
	adminUser := models.User{
		Name:     "Admin User",
		Email:    "admin@example.com",
		Password: string("password"),
		Role:     "admin",
	}
	
	result := db.Create(&adminUser)
	if result.Error != nil {
		log.Fatal("Failed to create admin user:", result.Error)
	}
	log.Printf("Admin user created: %s / password\n", adminUser.Email)
	
	regularUser := models.User{
		Name:     "Regular User",
		Email:    "user@example.com",
		Password: string("password"),
		Role:     "user",
	}
	
	result = db.Create(&regularUser)
	if result.Error != nil {
		log.Fatal("Failed to create regular user:", result.Error)
	}
	log.Printf("Regular user created: %s / password\n", regularUser.Email)
	
	sampleItems := []models.Item{
		{
			Name:        "Laptop Dell XPS 15",
			Description: "High-performance laptop with 16GB RAM, 512GB SSD",
			Stock:       10,
			Price:       25000000,
			Category:    "Electronics",
			CreatedBy:   adminUser.ID,
		},
		{
			Name:        "Office Desk",
			Description: "Wooden office desk 160x80 cm",
			Stock:       5,
			Price:       1500000,
			Category:    "Furniture",
			CreatedBy:   adminUser.ID,
		},
		{
			Name:        "Wireless Mouse",
			Description: "Logitech wireless mouse with USB receiver",
			Stock:       50,
			Price:       250000,
			Category:    "Accessories",
			CreatedBy:   adminUser.ID,
		},
		{
			Name:        "Office Chair",
			Description: "Ergonomic office chair with adjustable height",
			Stock:       8,
			Price:       3500000,
			Category:    "Furniture",
			CreatedBy:   adminUser.ID,
		},
		{
			Name:        "Monitor 24 inch",
			Description: "Full HD 24-inch monitor",
			Stock:       12,
			Price:       3000000,
			Category:    "Electronics",
			CreatedBy:   adminUser.ID,
		},
	}
	
	for i, item := range sampleItems {
		result = db.Create(&item)
		if result.Error != nil {
			log.Printf("Failed to create item %d: %v\n", i+1, result.Error)
			continue
		}
		log.Printf("Item created: %s (Stock: %d)\n", item.Name, item.Stock)
		
		activityLog := models.ActivityLog{
			UserID:      adminUser.ID,
			UserName:    adminUser.Name,
			ItemID:      item.ID,
			ItemName:    item.Name,
			Action:      "ITEM_CREATED",
			OldStock:    0,
			NewStock:    item.Stock,
			Description: "Item created during initial seeding",
		}
		
		if err := db.Create(&activityLog).Error; err != nil {
			log.Printf("Failed to create activity log for item %s: %v\n", item.Name, err)
		}
	}
	
	itemToUpdate := sampleItems[0]
	activityLog := models.ActivityLog{
		UserID:      adminUser.ID,
		UserName:    adminUser.Name,
		ItemID:      itemToUpdate.ID,
		ItemName:    itemToUpdate.Name,
		Action:      "STOCK_INCREMENT",
		Quantity:    5,
		OldStock:    itemToUpdate.Stock,
		NewStock:    itemToUpdate.Stock + 5,
		Description: "Stock incremented during initial seeding",
	}
	
	if err := db.Create(&activityLog).Error; err != nil {
		log.Printf("Failed to create stock update activity: %v\n", err)
	}
	
	activityLog2 := models.ActivityLog{
		UserID:      regularUser.ID,
		UserName:    regularUser.Name,
		ItemID:      sampleItems[1].ID,
		ItemName:    sampleItems[1].Name,
		Action:      "STOCK_DECREMENT",
		Quantity:    2,
		OldStock:    sampleItems[1].Stock,
		NewStock:    sampleItems[1].Stock - 2,
		Description: "Stock decremented during initial seeding",
	}
	
	if err := db.Create(&activityLog2).Error; err != nil {
		log.Printf("Failed to create stock update activity: %v\n", err)
	}
	
	log.Println("Sample data seeding completed successfully!")
	log.Println("==========================================")
	log.Println("Admin credentials:")
	log.Println("Email: admin@example.com")
	log.Println("Password: password")
	log.Println("")
	log.Println("User credentials:")
	log.Println("Email: user@example.com")
	log.Println("Password: password")
	log.Println("==========================================")
}
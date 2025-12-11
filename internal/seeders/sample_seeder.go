package seeders

import (
	"log"

	"inventory-api/internal/models"

	"gorm.io/gorm"
)

type SampleDataSeeder struct {
	DB *gorm.DB
}

func NewSampleDataSeeder(db *gorm.DB) *SampleDataSeeder {
	return &SampleDataSeeder{DB: db}
}

func (s *SampleDataSeeder) Run() error {
	log.Println("=== Starting sample data seeder ===")
	
	var adminUser models.User
	if err := s.DB.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		log.Printf("Warning: Admin user not found, skipping sample data seeding: %v", err)
		return nil
	}
	
	var regularUser models.User
	if err := s.DB.Where("email = ?", "user@example.com").First(&regularUser).Error; err != nil {
		if err := s.DB.Where("email = ?", "bob.johnson@example.com").First(&regularUser).Error; err != nil {
			log.Printf("Warning: Regular user not found, using admin for all activities")
			regularUser = adminUser
		}
	}
	
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
		result := s.DB.Where(models.Item{Name: item.Name}).FirstOrCreate(&item)
		if result.Error != nil {
			log.Printf("Failed to create item %d: %v\n", i+1, result.Error)
			continue
		}
		
		if result.RowsAffected > 0 {
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
			
			if err := s.DB.Create(&activityLog).Error; err != nil {
				log.Printf("Failed to create activity log for item %s: %v\n", item.Name, err)
			}
		} else {
			log.Printf("Item already exists: %s, skipping...\n", item.Name)
		}
	}
	
	var firstItem models.Item
	if err := s.DB.Where("name = ?", "Laptop Dell XPS 15").First(&firstItem).Error; err == nil {
		activityLog := models.ActivityLog{
			UserID:      adminUser.ID,
			UserName:    adminUser.Name,
			ItemID:      firstItem.ID,
			ItemName:    firstItem.Name,
			Action:      "STOCK_INCREMENT",
			Quantity:    5,
			OldStock:    firstItem.Stock,
			NewStock:    firstItem.Stock + 5,
			Description: "Stock incremented during initial seeding",
		}
		
		if err := s.DB.Create(&activityLog).Error; err != nil {
			log.Printf("Failed to create stock update activity: %v\n", err)
		}
	}
	
	var secondItem models.Item
	if err := s.DB.Where("name = ?", "Office Desk").First(&secondItem).Error; err == nil {
		activityLog2 := models.ActivityLog{
			UserID:      regularUser.ID,
			UserName:    regularUser.Name,
			ItemID:      secondItem.ID,
			ItemName:    secondItem.Name,
			Action:      "STOCK_DECREMENT",
			Quantity:    2,
			OldStock:    secondItem.Stock,
			NewStock:    secondItem.Stock - 2,
			Description: "Stock decremented during initial seeding",
		}
		
		if err := s.DB.Create(&activityLog2).Error; err != nil {
			log.Printf("Failed to create stock update activity: %v\n", err)
		}
	}
	
	log.Println("=== Sample data seeding completed! ===")
	log.Println("========================================")
	log.Println("Admin credentials:")
	log.Println("Email: admin@example.com")
	log.Println("Password: password")
	log.Println("")
	log.Println("User credentials:")
	log.Println("Email: user@example.com")
	log.Println("Password: password")
	log.Println("========================================")
	
	return nil
}
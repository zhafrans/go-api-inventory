package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"inventory-api/internal/config"
	"inventory-api/internal/controllers"
	"inventory-api/internal/database"
	"inventory-api/internal/middleware"
	"inventory-api/internal/seeders"
	"inventory-api/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
	
	cfg := config.LoadConfig()
	
	database.ConnectDB(cfg)
	
	runSeeders()
	
	services.NewItemService()
	
	authController := controllers.NewAuthController(cfg)
	itemController := controllers.NewItemController()
	activityController := controllers.NewActivityController()
	
	app := fiber.New(fiber.Config{
		AppName: "Inventory Management API",
	})
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE",
	}))
	app.Use(logger.New())
	
	api := app.Group("/api")
	api.Post("/register", authController.Register)
	api.Post("/login", authController.Login)
	
	protected := api.Group("", middleware.JWTMiddleware(cfg))
	protected.Get("/profile", authController.Profile)

	protected.Get("/activities", activityController.GetAllActivities)
	
	items := protected.Group("/items")
	items.Post("/", itemController.CreateItem)
	items.Get("/", itemController.GetAllItems)
	items.Get("/:id", itemController.GetItemByID)
	items.Put("/:id", itemController.UpdateItem)
	items.Patch("/:id/stock", itemController.UpdateStock)
	items.Delete("/:id", itemController.DeleteItem)

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := app.Listen(cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func runSeeders() {
	log.Println("=== Running Database Seeders ===")
	
	if database.DB == nil {
		log.Fatal("Database connection is not initialized!")
	}
	
	sampleSeeder := seeders.NewSampleDataSeeder(database.DB)
	if err := sampleSeeder.Run(); err != nil {
		log.Printf("Warning: Sample data seeder failed: %v", err)
	}
	
	log.Println("=== All seeders completed ===")
}
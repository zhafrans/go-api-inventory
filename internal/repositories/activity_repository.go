package repositories

import (
	"inventory-api/internal/database"
	"inventory-api/internal/models"

	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository() *ActivityRepository {
	return &ActivityRepository{db: database.DB}
}

func (r *ActivityRepository) Create(activity *models.ActivityLog) error {
	return database.DB.Create(activity).Error
}

func (r *ActivityRepository) FindByItemID(itemID string) ([]models.ActivityLog, error) {
	var activities []models.ActivityLog
	err := database.DB.Where("item_id = ?", itemID).
		Order("created_at DESC").
		Find(&activities).Error
	return activities, err
}

func (r *ActivityRepository) FindAll() ([]models.ActivityLog, error) {
	var activities []models.ActivityLog
	err := database.DB.Order("created_at DESC").Find(&activities).Error
	return activities, err
}
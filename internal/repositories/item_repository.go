package repositories

import (
	"inventory-api/internal/database"
	"inventory-api/internal/models"

	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository() *ItemRepository {
	return &ItemRepository{db: database.DB}
}

func (r *ItemRepository) Create(item *models.Item) error {
	return database.DB.Create(item).Error
}

func (r *ItemRepository) FindAll() ([]models.Item, error) {
	var items []models.Item
	err := r.db.Preload("Creator", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Order("created_at DESC").Find(&items).Error
	return items, err
}


func (r *ItemRepository) FindByID(id string) (*models.Item, error) {
	var item models.Item
	err := r.db.Preload("Creator", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Where("id = ?", id).First(&item).Error
	return &item, err
}

func (r *ItemRepository) Update(item *models.Item) error {
	return database.DB.Save(item).Error
}

func (r *ItemRepository) Delete(id string) error {
	return database.DB.Where("id = ?", id).Delete(&models.Item{}).Error
}

func (r *ItemRepository) UpdateStock(itemID string, quantity int) error {
	return database.DB.Model(&models.Item{}).
		Where("id = ?", itemID).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
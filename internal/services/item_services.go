package services

import (
	"errors"

	"inventory-api/internal/models"
	"inventory-api/internal/repositories"
)

type ItemService struct {
	itemRepo *repositories.ItemRepository
	activityRepo *repositories.ActivityRepository
	userRepo *repositories.UserRepository
}

func NewItemService() *ItemService {
	return &ItemService{
		itemRepo:     repositories.NewItemRepository(),
		activityRepo: repositories.NewActivityRepository(),
		userRepo:     repositories.NewUserRepository(),
	}
}

func (s *ItemService) CreateItem(req *models.CreateItemRequest, userID string) (*models.Item, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	item := &models.Item{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Stock:       req.Stock,
		MinStock:    req.MinStock,
		MaxStock:    req.MaxStock,
		Price:       req.Price,
		SKU:         req.SKU,
		Location:    req.Location,
		CreatedBy:   userID,
	}
	
	if err := s.itemRepo.Create(item); err != nil {
		return nil, err
	}
	
	activity := &models.ActivityLog{
		UserID:       userID,
		UserName:     user.Name,
		ItemID:       item.ID,
		ItemName:     item.Name,
		Action: models.ActivityTypeItemCreated,
		Quantity:     req.Stock,
		OldStock:     0,
		NewStock:     req.Stock,
		Description:  "Item created",
	}
	s.activityRepo.Create(activity)
	
	return item, nil
}

func (s *ItemService) GetAllItems() ([]models.Item, error) {
	return s.itemRepo.FindAll()
}

func (s *ItemService) GetItemByID(id string) (*models.Item, error) {
	return s.itemRepo.FindByID(id)
}

func (s *ItemService) UpdateItem(id string, req *models.UpdateItemRequest, userID string) (*models.Item, error) {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("item not found")
	}
	
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.MinStock > 0 {
		item.MinStock = req.MinStock
	}
	if req.MaxStock > 0 {
		item.MaxStock = req.MaxStock
	}
	if req.Price > 0 {
		item.Price = req.Price
	}
	if req.Location != "" {
		item.Location = req.Location
	}
	
	if err := s.itemRepo.Update(item); err != nil {
		return nil, err
	}
	
	activity := &models.ActivityLog{
		UserID:       userID,
		UserName:     user.Name,
		ItemID:       item.ID,
		ItemName:     item.Name,
		Action: models.ActivityTypeItemUpdated,
		Description:  "Item updated",
	}
	s.activityRepo.Create(activity)
	
	return item, nil
}

func (s *ItemService) UpdateStock(id string, req *models.UpdateStockRequest, userID string) (*models.Item, error) {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("item not found")
	}
	
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	oldStock := item.Stock
	var newStock int
	var activityType models.ActivityType
	
	if req.Type == "increment" {
		newStock = item.Stock + req.Quantity
		activityType = models.ActivityTypeStockIncrement
	} else {
		newStock = item.Stock - req.Quantity
		if newStock < 0 {
			return nil, errors.New("insufficient stock")
		}
		activityType = models.ActivityTypeStockDecrement
	}
	
	item.Stock = newStock
	if err := s.itemRepo.Update(item); err != nil {
		return nil, err
	}
	
	activity := &models.ActivityLog{
		UserID:       userID,
		UserName:     user.Name,
		ItemID:       item.ID,
		ItemName:     item.Name,
		Action: activityType,
		Quantity:     req.Quantity,
		OldStock:     oldStock,
		NewStock:     newStock,
		Description:  req.Reason,
	}
	s.activityRepo.Create(activity)
	
	return item, nil
}

func (s *ItemService) DeleteItem(id string, userID string) error {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return errors.New("item not found")
	}
	
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	
	if err := s.itemRepo.Delete(id); err != nil {
		return err
	}
	
	activity := &models.ActivityLog{
		UserID:       userID,
		UserName:     user.Name,
		ItemID:       item.ID,
		ItemName:     item.Name,
		Action: models.ActivityTypeItemDeleted,
		Description:  "Item deleted",
	}
	s.activityRepo.Create(activity)
	
	return nil
}
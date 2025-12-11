package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"inventory-api/internal/models"
	"inventory-api/internal/services"
)

type ItemController struct {
	itemService     *services.ItemService
	activityService *services.ActivityService
	responseService *services.ResponseService
}

func NewItemController() *ItemController {
	return &ItemController{
		itemService:     services.NewItemService(),
		activityService: services.NewActivityService(),
		responseService: services.NewResponseService(),
	}
}

func (ctrl *ItemController) CreateItem(c *fiber.Ctx) error {
	var req models.CreateItemRequest
	
	if err := c.BodyParser(&req); err != nil {
		return ctrl.responseService.BadRequest(c, "Invalid request body", err.Error())
	}
	
	if req.Name == "" {
		return ctrl.responseService.BadRequest(c, "Validation failed", "Item name is required")
	}

	fmt.Println("Item name:", req.Name)
	
	userID := c.Locals("userID")
	if userID == nil {
		return ctrl.responseService.Unauthorized(c, "Authentication required", "User not authenticated")
	}
	
	userIDStr, ok := userID.(string)
	if !ok {
		return ctrl.responseService.Unauthorized(c, "Invalid user session", "Invalid user ID format")
	}
	
	_, err := ctrl.itemService.CreateItem(&req, userIDStr)
	if err != nil {
		return ctrl.responseService.BadRequest(c, "Failed to create item", err.Error())
	}
	
	return ctrl.responseService.Created(c, "Item created successfully", nil)
}

func (ctrl *ItemController) GetAllItems(c *fiber.Ctx) error {
	items, err := ctrl.itemService.GetAllItems()
	if err != nil {
		return ctrl.responseService.InternalServerError(c, "Failed to fetch items", err.Error())
	}
	
	return ctrl.responseService.Success(c, fiber.StatusOK, "Items retrieved successfully", fiber.Map{
		"items": items,
		"count": len(items),
	})
}

func (ctrl *ItemController) GetItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	
	item, err := ctrl.itemService.GetItemByID(id)
	if err != nil {
		return ctrl.responseService.NotFound(c, "Item not found", err.Error())
	}
	
	return ctrl.responseService.Success(c, fiber.StatusOK, "Item retrieved successfully", fiber.Map{
		"item": item,
	})
}

func (ctrl *ItemController) UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	
	currentItem, err := ctrl.itemService.GetItemByID(id)
	if err != nil {
		return ctrl.responseService.NotFound(c, "Item not found", err.Error())
	}
	
	var req models.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return ctrl.responseService.BadRequest(c, "Invalid request body", err.Error())
	}
	
	userID := c.Locals("userID")
	if userID == nil {
		return ctrl.responseService.Unauthorized(c, "Authentication required", "User not authenticated")
	}
	
	userIDStr, ok := userID.(string)
	if !ok {
		return ctrl.responseService.Unauthorized(c, "Invalid user session", "Invalid user ID format")
	}
	
	updatedItem, err := ctrl.itemService.UpdateItem(id, &req, userIDStr)
	if err != nil {
		return ctrl.responseService.BadRequest(c, "Failed to update item", err.Error())
	}
	
	changes := fiber.Map{}
	if currentItem.Name != updatedItem.Name {
		changes["name"] = fiber.Map{
			"from": currentItem.Name,
			"to":   updatedItem.Name,
		}
	}
	if currentItem.Description != updatedItem.Description {
		changes["description"] = fiber.Map{
			"from": currentItem.Description,
			"to":   updatedItem.Description,
		}
	}
	if currentItem.Price != updatedItem.Price {
		changes["price"] = fiber.Map{
			"from": currentItem.Price,
			"to":   updatedItem.Price,
		}
	}
	if currentItem.Category != updatedItem.Category {
		changes["category"] = fiber.Map{
			"from": currentItem.Category,
			"to":   updatedItem.Category,
		}
	}
	
	return ctrl.responseService.Success(c, fiber.StatusOK, "Item updated successfully", fiber.Map{
		"item":    updatedItem,
		"changes": changes,
	})
}

func (ctrl *ItemController) UpdateStock(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var req models.UpdateStockRequest
	if err := c.BodyParser(&req); err != nil {
		return ctrl.responseService.BadRequest(c, "Invalid request body", err.Error())
	}
	
	if req.Quantity <= 0 {
		return ctrl.responseService.BadRequest(c, "Validation failed", "Quantity must be greater than 0")
	}
	
	if req.Type != "increment" && req.Type != "decrement" {
		return ctrl.responseService.BadRequest(c, "Validation failed", "Type must be 'increment' or 'decrement'")
	}
	
	userID := c.Locals("userID")
	if userID == nil {
		return ctrl.responseService.Unauthorized(c, "Authentication required", "User not authenticated")
	}
	
	userIDStr, ok := userID.(string)
	if !ok {
		return ctrl.responseService.Unauthorized(c, "Invalid user session", "Invalid user ID format")
	}
	
	updatedItem, err := ctrl.itemService.UpdateStock(id, &req, userIDStr)
	if err != nil {
		return ctrl.responseService.BadRequest(c, "Failed to update stock", err.Error())
	}
	
	return ctrl.responseService.Success(c, fiber.StatusOK, "Stock updated successfully", fiber.Map{
		"item": updatedItem,
	})
}

func (ctrl *ItemController) DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")
	
	item, err := ctrl.itemService.GetItemByID(id)
	if err != nil {
		return ctrl.responseService.NotFound(c, "Item not found", err.Error())
	}
	
	userID := c.Locals("userID")
	if userID == nil {
		return ctrl.responseService.Unauthorized(c, "Authentication required", "User not authenticated")
	}
	
	userIDStr, ok := userID.(string)
	if !ok {
		return ctrl.responseService.Unauthorized(c, "Invalid user session", "Invalid user ID format")
	}
	
	err = ctrl.itemService.DeleteItem(id, userIDStr)
	if err != nil {
		return ctrl.responseService.BadRequest(c, "Failed to delete item", err.Error())
	}
	
	return ctrl.responseService.Success(c, fiber.StatusOK, "Item deleted successfully", fiber.Map{
		"deleted_item": fiber.Map{
			"id":   item.ID,
			"name": item.Name,
		},
	})
}
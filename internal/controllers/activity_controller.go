package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"inventory-api/internal/services"
)

type ActivityController struct {
	activityService *services.ActivityService
	responseService *services.ResponseService
}

func NewActivityController() *ActivityController {
	return &ActivityController{
		activityService: services.NewActivityService(),
		responseService: services.NewResponseService(),
	}
}

func (ctrl *ActivityController) GetAllActivities(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	activityType := c.Query("type", "")
	itemID := c.Query("item_id", "")
	userID := c.Query("user_id", "")
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	activities, total, err := ctrl.activityService.GetAllActivities(page, limit, activityType, itemID, userID)
	if err != nil {
		return ctrl.responseService.InternalServerError(c, "Failed to fetch activities", err.Error())
	}

	return ctrl.responseService.SuccessWithPagination(
		c,
		fiber.StatusOK,
		"Activities retrieved successfully",
		activities,
		page,
		limit,
		total,
	)
}

func (ctrl *ActivityController) GetItemActivities(c *fiber.Ctx) error {
	itemID := c.Params("id")
	
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	activities, total, err := ctrl.activityService.GetActivitiesByItemID(itemID, page, limit)
	if err != nil {
		return ctrl.responseService.InternalServerError(c, "Failed to fetch item activities", err.Error())
	}

	return ctrl.responseService.SuccessWithPagination(
		c,
		fiber.StatusOK,
		"Item activities retrieved successfully",
		activities,
		page,
		limit,
		total,
	)
}

func (ctrl *ActivityController) GetRecentActivities(c *fiber.Ctx) error {
	activities, err := ctrl.activityService.GetRecentActivities(10)
	if err != nil {
		return ctrl.responseService.InternalServerError(c, "Failed to fetch recent activities", err.Error())
	}

	return ctrl.responseService.Success(c, fiber.StatusOK, "Recent activities retrieved successfully", fiber.Map{
		"activities": activities,
		"count":      len(activities),
	})
}
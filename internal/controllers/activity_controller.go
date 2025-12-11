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
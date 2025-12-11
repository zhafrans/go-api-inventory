package controllers

import (
	"github.com/gofiber/fiber/v2"

	"inventory-api/internal/config"
	"inventory-api/internal/models"
	"inventory-api/internal/services"
)

type AuthController struct {
	authService    *services.AuthService
	config         *config.Config
	responseService *services.ResponseService
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		authService:    services.NewAuthService(cfg),
		config:         cfg,
		responseService: services.NewResponseService(),
	}
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	
	if err := c.BodyParser(&req); err != nil {
		return ctrl.responseService.BadRequest(c, "Invalid request body", err.Error())
	}
	
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return ctrl.responseService.BadRequest(c, "Validation failed", "Name, email, and password are required")
	}
	
	user, err := ctrl.authService.Register(&req)
	if err != nil {
		return ctrl.responseService.BadRequest(c, "Registration failed", err.Error())
	}
	
	return ctrl.responseService.Created(c, "Registration successful", fiber.Map{
		"user": user,
	})
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	
	if err := c.BodyParser(&req); err != nil {
		return ctrl.responseService.BadRequest(c, "Invalid request body", err.Error())
	}
	
	if req.Email == "" || req.Password == "" {
		return ctrl.responseService.BadRequest(c, "Validation failed", "Email and password are required")
	}
	
	token, err := ctrl.authService.Login(&req)
	if err != nil {
		return ctrl.responseService.Unauthorized(c, "Login failed", err.Error())
	}
	
	return ctrl.responseService.Success(c, fiber.StatusOK, "Login successful", fiber.Map{
		"token": token,
	})
}

func (ctrl *AuthController) Profile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return ctrl.responseService.Unauthorized(c, "Authentication required", "Invalid user session")
	}
	
	user, err := ctrl.authService.GetUserProfile(userID)
	if err != nil {
		return ctrl.responseService.NotFound(c, "User not found", err.Error())
	}
	
	return ctrl.responseService.Success(c, fiber.StatusOK, "Profile retrieved successfully", fiber.Map{
		"user": user,
	})
}
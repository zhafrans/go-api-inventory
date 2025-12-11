package services

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type ResponseService struct{}

func NewResponseService() *ResponseService {
	return &ResponseService{}
}

func (rs *ResponseService) Success(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(Response{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func (rs *ResponseService) SuccessWithMeta(c *fiber.Ctx, code int, message string, data interface{}, meta interface{}) error {
	return c.Status(code).JSON(Response{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func (rs *ResponseService) Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:  "success",
		Code:    fiber.StatusCreated,
		Message: message,
		Data:    data,
	})
}

func (rs *ResponseService) Error(c *fiber.Ctx, code int, message string, errDetail interface{}) error {
	return c.Status(code).JSON(Response{
		Status:  "error",
		Code:    code,
		Message: message,
		Error:   errDetail,
	})
}

func (rs *ResponseService) BadRequest(c *fiber.Ctx, message string, errDetail interface{}) error {
	return rs.Error(c, fiber.StatusBadRequest, message, errDetail)
}

func (rs *ResponseService) Unauthorized(c *fiber.Ctx, message string, errDetail interface{}) error {
	return rs.Error(c, fiber.StatusUnauthorized, message, errDetail)
}

func (rs *ResponseService) Forbidden(c *fiber.Ctx, message string, errDetail interface{}) error {
	return rs.Error(c, fiber.StatusForbidden, message, errDetail)
}

func (rs *ResponseService) NotFound(c *fiber.Ctx, message string, errDetail interface{}) error {
	return rs.Error(c, fiber.StatusNotFound, message, errDetail)
}

func (rs *ResponseService) InternalServerError(c *fiber.Ctx, message string, errDetail interface{}) error {
	return rs.Error(c, fiber.StatusInternalServerError, message, errDetail)
}

func (rs *ResponseService) ValidationError(c *fiber.Ctx, message string, validationErrors interface{}) error {
	return rs.Error(c, fiber.StatusBadRequest, message, validationErrors)
}
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func (rs *ResponseService) SuccessWithPagination(
	c *fiber.Ctx,
	code int,
	message string,
	data interface{},
	page int,
	limit int,
	total int64,
) error {
	totalPages := (total + int64(limit) - 1) / int64(limit)
	
	meta := PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
		HasNext:    int64(page) < totalPages,
		HasPrev:    page > 1,
	}
	
	return rs.SuccessWithMeta(c, code, message, data, meta)
}
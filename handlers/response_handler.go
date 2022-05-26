package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ResponseParams struct {
	Data    *fiber.Map
	Message string
}

func SendSuccessResponse(c *fiber.Ctx, data *fiber.Map) error {
	return c.Status(http.StatusAccepted).JSON(data)
}

func SendBadRequestResponse(c *fiber.Ctx, data *fiber.Map) error {
	return c.Status(http.StatusBadRequest).JSON(data)
}

func SendErrorResponse(c *fiber.Ctx, data *fiber.Map) error {
	return c.Status(http.StatusInternalServerError).JSON(data)
}

func SendBadAuthResponse(c *fiber.Ctx, data *fiber.Map) error {
	return c.Status(http.StatusUnauthorized).JSON(data)
}

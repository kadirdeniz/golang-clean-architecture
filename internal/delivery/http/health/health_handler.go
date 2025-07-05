package health

import (
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Health(c *fiber.Ctx) error
}

type handler struct {}

func NewHandler() Handler {
	return &handler{}
}

// @Summary Health check endpoint
// @Description Check if the application is healthy
// @Tags health
// @Produce json
// @Success 200 
// @Router /health [get]
func (h *handler) Health(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

package health

import "github.com/gofiber/fiber/v2"

type Router interface {
	RegisterRoutes(app *fiber.App)
}

type router struct {
	handler Handler
}

func NewRouter(handler Handler) Router {
	return &router{
		handler: handler,
	}
}

func (r *router) RegisterRoutes(app *fiber.App) {
	app.Get("/health", r.handler.Health)
}
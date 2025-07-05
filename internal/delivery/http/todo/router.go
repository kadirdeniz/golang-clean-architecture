package todo

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
	todos := app.Group("/todos")
	todos.Post("/", r.handler.Create)
	todos.Get("/", r.handler.GetAll)
	todos.Get("/:id", r.handler.GetByID)
	todos.Put("/:id", r.handler.Update)
	todos.Delete("/:id", r.handler.Delete)
}
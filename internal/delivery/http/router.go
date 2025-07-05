package http

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"

	"github.com/kadirdeniz/golang-clean-architecture/internal/delivery/http/health"
	"github.com/kadirdeniz/golang-clean-architecture/internal/delivery/http/todo"
)

type Router interface {
	App() *fiber.App
}

type router struct {
	app *fiber.App
}

func NewRouter(todoRouter todo.Router,healthRouter health.Router) Router {
	app := fiber.New()

	app.Use(limiter.New())	
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error} | ${pid} ${locals:requestid}\n",
		TimeFormat: "15:04:05",
		TimeZone:   "UTC",
		Output:     os.Stdout,

	}))

	app.Use(cors.New(cors.Config{           
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	healthRouter.RegisterRoutes(app)
	todoRouter.RegisterRoutes(app)

	return &router{
		app: app,
	}
}	

func (r *router) App() *fiber.App {
	return r.app
}

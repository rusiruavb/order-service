package routes

import (
	"order_service/controllers"

	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(app *fiber.App) {
	app.Post("/order", controllers.CreateOrder)
	app.Get("/order/user", controllers.GetOrders)
	app.Get("/order/:orderId", controllers.GetOrderById)
}

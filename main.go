package main

import (
	"fmt"
	"order_service/configs"
	"order_service/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("📦 Order Service v1.0.1")
	})

	routes.OrderRoutes(app)

	configs.ConnectDB()
	fmt.Println("✨ Database Synced")
	fmt.Println("📦 Order service is running")
	app.Listen(":9090")
}

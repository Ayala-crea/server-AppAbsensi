package main

import (
	"Ayala-Crea/server/config"
	"Ayala-Crea/server/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	client := config.ConnectDb()
	config.InitFirebase()

	app.Use(logger.New(logger.Config{
		Format: "${status} - ${method} ${path}\n",
	}))

	app.Use(cors.New(cors.Config{ 
		AllowHeaders: "*",
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", client.Database("appAbsensi"))
		return c.Next()
	})

	routes.SetupTaskRoutes(app)

	// Menambahkan rute untuk menguji apakah server berjalan
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running. Go to /img/<filename> to view the image.")
	})

	// Menjalankan server Fiber pada port 3000
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
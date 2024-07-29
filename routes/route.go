package routes

import (
	"Ayala-Crea/server-app-absensi/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupTaskRoutes(app *fiber.App) {
	app.Post("/register", controllers.RegisterUser)
}
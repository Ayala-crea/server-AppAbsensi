package routes

import (
	"Ayala-Crea/server/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupTaskRoutes(app *fiber.App) {
	app.Post("/register", controllers.RegisterUser)
	app.Post("/login", controllers.LoginAdmin)
	app.Post("/profile", controllers.CreateProfile)
	app.Get("/profile/me", controllers.GetProfileById)
	app.Post("/dataAbsen", controllers.InsertDataEmployeeOrStudents)
}
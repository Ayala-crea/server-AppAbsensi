package controllers

import (
	"Ayala-Crea/server-app-absensi/models"
	repo "Ayala-Crea/server-app-absensi/repository"
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	var user models.Users

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "request Body Invalid",
		})
	}

	db := c.Locals("db").(*sql.DB)

    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&count)
    if err != nil {
        log.Println("Error querying email:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal cek email",
        })
    }
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
            "message": "Email sudah terdaftar",
        })
	}

	if !user.IsActive {
		user.IsActive = true
	}

	user.CreatedAt = time.Now()

	// Menyimpan user ke database menggunakan repository
    if err := repo.CreateUser(db, &user); err != nil {
        log.Println(err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal Register",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Berhasil Register!",
		"data": user,
    })
}
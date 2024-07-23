package controllers

import (
	"Ayala-Crea/server/models"
	repo "Ayala-Crea/server/repository"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *fiber.Ctx) error {
	var user models.Users

	// Parsing body request ke struct User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request Body Invalid",
		})
	}

 	// Mengambil instance MongoDB dari context
	db := c.Locals("db").(*mongo.Database)

	// Cek apakah email sudah digunakan
	filter := bson.M{"email": user.Email}
	count, err := db.Collection("users").CountDocuments(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal cek email",
		})
	}
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email sudah terdaftar",
		})
	}

	// Atur nilai default untuk isActive jika kosong
	if user.IsActive == "" {
		user.IsActive = "active"
	}

	// Menyimpan user ke database menggunakan repository
	if err := repo.CreateUser(db, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal Register",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil Register!",
	})
}

func LoginAdmin(c *fiber.Ctx) error {
    var user models.Users
    db := c.Locals("db").(*mongo.Database)

    // Parsing request body into user struct
    if err := c.BodyParser(&user); err != nil {
        fmt.Println("Error parsing request data:", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Error parsing request data",
        })
    }
    fmt.Println("Parsed user:", user)

    // Validate email or username provided
    if user.Email == "" && user.Username == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Request must include email or username",
        })
    }

    // Retrieve user data from the database
    var userData *models.Users
    var err error

    if user.Email != "" {
        userData, err = repo.GetAdminByEmail(db, user.Email)
    } else if user.Username != "" {
        userData, err = repo.GetAdminByUsername(db, user.Username)
    }

    if err != nil || userData == nil {
        fmt.Println("Error retrieving user data or user not found:", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Username atau password salah",
        })
    }
    fmt.Println("Retrieved user data:", userData)

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password)); err != nil {
        fmt.Println("Error verifying password:", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Username atau password salah",
        })
    }

    // Generate JWT token
    token, err := repo.GenerateToken(userData)
    if err != nil {
        fmt.Println("Error generating token:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to generate token",
        })
    }

    // Get detailed user data by ID
    detailedUserData, err := repo.GetAdminById(db, userData.IdUser)
    if err != nil {
        fmt.Println("Error retrieving detailed user data:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to retrieve user data",
        })
    }

    // Return successful response with token and user data
    return c.JSON(fiber.Map{
        "token": token,
        "user": fiber.Map{
            "id_user": detailedUserData.IdUser,
            "id_role": detailedUserData.IdRole,
            "nama": detailedUserData.Nama,
            "username": detailedUserData.Username,
            "email": detailedUserData.Email,
            // Add more user details here as needed
        },
    })
}

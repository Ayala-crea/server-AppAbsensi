package controllers

import (
	"Ayala-Crea/server/models"
	repo "Ayala-Crea/server/repository"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)


func CreateProfile(c *fiber.Ctx) error {

	tokenStr := c.Get("login")
	if tokenStr == "" {
		log.Println("Error: Header 'login' is missing.")
		return fiber.NewError(fiber.StatusBadRequest, "Header tidak ada")
	}

	// Parse token untuk mendapatkan id
	token, err := jwt.ParseWithClaims(tokenStr, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil // Ganti "secret_key" dengan kunci rahasia Anda
	})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid"})
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok || !token.Valid {
		log.Println("Error: Invalid token.")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid"})
	}

	// ID pengguna yang diambil dari token
	idUser := claims.IdUser

	var profile models.Profile

	if err := c.BodyParser(&profile); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	file, err := c.FormFile("img")
	if err != nil {
		log.Printf("Error: Image not found in form data: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image Not Found"})
	}

	// Baca file
	fileContent, err := file.Open()
	if err != nil {
		log.Printf("Error reading image file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not read image file"})
	}
	defer fileContent.Close()

	// Membaca file ke dalam byte slice
	fileBytes := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBytes, fileContent); err != nil {
		log.Printf("Error copying file content: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not read image file"})
	}

	imageURL, err := repo.UploadToFirebaseStorage(fileBytes.Bytes(), file.Filename)
	if err != nil {
		log.Printf("Error uploading image to Firebase: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Could not upload image to Firebase: %v", err)})
	}

	profile.Image = imageURL
	profile.IdUser = idUser // Simpan ID pengguna langsung dari klaim token

	db := c.Locals("db").(*mongo.Database)
	if err := repo.InsertProfile(db, &profile); err != nil {
		log.Printf("Error saving profile to database: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could save data profile"})
	}

	log.Println("Profile created successfully.")
	return c.Status(http.StatusCreated).JSON(fiber.Map{"code": http.StatusCreated, "success": true, "status": "success", "message": "Task berhasil disimpan", "data": profile})
}
func GetProfileById(c *fiber.Ctx) error {
	// Cek token header autentikasi
	token := c.Get("login")
	if token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Header tidak ada")
	}

	id := c.Query("profile_id")
	if id == "" {
		return fiber.NewError(http.StatusBadRequest, "id tidak ditemukan")
	}

	db := c.Locals("db").(*mongo.Database)

	profile, err := repo.GetProfileById(db, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"code": http.StatusOK, "success": true, "status": "success", "data": profile})
}

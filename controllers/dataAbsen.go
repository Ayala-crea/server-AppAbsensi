package controllers

import (
	"log"

	"Ayala-Crea/server/models"
	repo "Ayala-Crea/server/repository"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertDataEmployeeOrStudents(c *fiber.Ctx) error {
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

	// Read the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File tidak ada"})
	}

	// Open the uploaded file
	f, err := file.Open()
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error opening file"})
	}
	defer f.Close()

	// Read the Excel file
	xlFile, err := excelize.OpenReader(f)
	if err != nil {
		log.Printf("Error reading Excel file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error reading Excel file"})
	}

	// Iterate over the rows and insert data
	rows, err := xlFile.GetRows("Sheet1")
	if err != nil {
		log.Printf("Error getting rows from Excel file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error getting rows from Excel file"})
	}

	db := c.Locals("db").(*mongo.Database)

	for _, row := range rows[1:] { // Skipping header row
		dataAbsen := models.Students_Employees{
			IdUser:      idUser,
			FullName:    row[1],
			Email:       row[2],
			Status:      row[3],
			Class:       row[4],
			NpkOrNpm:    row[5],
			PhoneNumber: row[6],
		}
		err = repo.InsertEmployeeOrStudent(db, &dataAbsen)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error inserting data"})
		}

		user := models.Users{
			IdRole:   2,
			Nama:     row[1],
			Username: row[4],
			Password: "User123",
			Email:    row[2],
			PhoneNumber: row[6],
			IsActive: "Active",
		}
		err := repo.CreateUser(db, &user)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error inserting user"})
		}
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Data inserted successfully"})
}

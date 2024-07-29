package repository

import (
	"Ayala-Crea/server-app-absensi/models"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, user *models.Users) error {
    // Hash the user's password
    hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashPassword)

    // SQL statement to insert user into the database
    sqlStatement := `
        INSERT INTO users (id_role, id_penginputan, nama, username, password, email, phone_number, created_at, is_active)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id`
    
    // Execute the SQL statement
    err = db.QueryRow(sqlStatement, user.IdRole, user.IdPenginputan, user.Nama, user.Username, user.Password, user.Email, user.PhoneNumber, user.CreatedAt, user.IsActive).Scan(&user.ID)
    if err != nil {
        return err
    }

    return nil
}
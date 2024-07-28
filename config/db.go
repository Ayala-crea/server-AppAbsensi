package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file
func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

// ConnectDb connects to the PostgreSQL database and returns the connection
func ConnectDb() *sql.DB {
    LoadEnv()

    pgHost := os.Getenv("PG_HOST")
    pgPort := os.Getenv("PG_PORT")
    pgUser := os.Getenv("PG_USER")
    pgPassword := os.Getenv("PG_PASSWORD")
    pgDBName := os.Getenv("PG_DBNAME")

    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        pgHost, pgPort, pgUser, pgPassword, pgDBName)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal("Failed to open a DB connection: ", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Failed to connect to the database: ", err)
    }

    fmt.Println("Connected to PostgreSQL!")

    return db
}

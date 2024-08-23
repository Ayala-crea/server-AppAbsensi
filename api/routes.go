package api

import (
	"database/sql"

	handlers "Ayala-Crea/server-app-absensi/api/handlers"

	"github.com/gorilla/mux"
)

func AllRoutes(r *mux.Router, db *sql.DB) {
	// Register the POST route for user registration
	r.HandleFunc("/register", handlers.RegisterUser(db)).Methods("POST")
	r.HandleFunc("/login", handlers.Login(db)).Methods("POST")

	r.HandleFunc("/upload", handlers.UploadExcel(db)).Methods("POST")
	r.HandleFunc("/data", handlers.GetAllStudentsEmployees(db)).Methods("GET")
}
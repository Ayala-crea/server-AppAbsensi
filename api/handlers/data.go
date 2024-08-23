package handlers

import (
	"Ayala-Crea/server-app-absensi/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"
)

func UploadExcel(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// 2. Validasi token JWT
		tokenString := strings.Split(authHeader, "Bearer ")[1]
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// id_penginputan diambil dari JWT
		idPenginputan := claims.IDPenginputan
		fmt.Printf("Validated id_penginputan: %d\n", idPenginputan)

		// AdminID diambil dari JWT
		adminID := claims.ID // Asumsikan ID di JWT adalah AdminID
		fmt.Printf("Validated AdminID: %d\n", adminID)

		// 3. Parse multipart form untuk mengunggah file
		err = r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// 4. Retrieve file dari form data
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 5. Buka file Excel menggunakan excelize
		f, err := excelize.OpenReader(file)
		if err != nil {
			http.Error(w, "Unable to read Excel file", http.StatusBadRequest)
			return
		}

		// 6. Dapatkan baris dari sheet pertama
		rows, err := f.GetRows(f.GetSheetName(0))
		if err != nil {
			http.Error(w, "Unable to read rows from Excel file", http.StatusBadRequest)
			return
		}

		// 7. Iterate melalui setiap baris dan simpan ke database
		for i, row := range rows {
			if i == 0 {
				// Lewatkan baris pertama jika itu adalah header
				continue
			}

			// Gunakan AdminID dari JWT
			student := models.StudentsEmployees{
				AdminID:     adminID, // Mengambil AdminID dari JWT token
				FullName:    row[1],
				Status:      row[2],
				Class:       row[3],
				NpkOrNpm:    row[4],
				PhoneNumber: row[5],
			}

			query := `INSERT INTO students_employees (admin_id, full_name, status, class, npk_or_npm, phone_number)
                      VALUES ($1, $2, $3, $4, $5, $6)`

			_, err := db.Exec(query, student.AdminID, student.FullName, student.Status, student.Class, student.NpkOrNpm, student.PhoneNumber)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to insert record for row %d: %v", i+1, err), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func convertToInt(value string) int {
	var result int
	fmt.Sscanf(value, "%d", &result)
	return result
}

func GetAllStudentsEmployees(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// 2. Validasi token JWT
		tokenString := strings.Split(authHeader, "Bearer ")[1]
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Query untuk mengambil semua data dari tabel students_employees
		query := `SELECT id, admin_id, full_name, status, class, npk_or_npm, phone_number FROM students_employees`

		// Eksekusi query
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Failed to execute query", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Slice untuk menampung hasil query
		var studentsEmployees []models.StudentsEmployees

		// Iterasi melalui hasil query
		for rows.Next() {
			var student models.StudentsEmployees
			err := rows.Scan(&student.ID, &student.AdminID, &student.FullName, &student.Status, &student.Class, &student.NpkOrNpm, &student.PhoneNumber)
			if err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			studentsEmployees = append(studentsEmployees, student)
		}

		// Cek error setelah iterasi selesai
		if err = rows.Err(); err != nil {
			http.Error(w, "Error occurred during iteration", http.StatusInternalServerError)
			return
		}

		// Kembalikan hasil dalam bentuk JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(studentsEmployees)
	}
}

func GetDataByIdAdmin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// 2. Validasi token JWT
		tokenString := strings.Split(authHeader, "Bearer ")[1]
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		params := mux.Vars(r)
        idAdmin := convertToInt(params["id_admin"])
		query := `SELECT id, admin_id, full_name, status, class, npk_or_npm, phone_number FROM students_employees WHERE admin_id = $1`
		row := db.QueryRow(query, idAdmin)
		var student models.StudentsEmployees
		err = row.Scan(&student.ID, &student.AdminID, &student.FullName, &student.Status, &student.Class, &student.NpkOrNpm, &student.PhoneNumber)
		if err == sql.ErrNoRows {
			http.Error(w, "Record not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to execute query", http.StatusInternalServerError)
			return
		}
	}
}
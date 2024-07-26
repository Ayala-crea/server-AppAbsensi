package repository

import (
	"Ayala-Crea/server/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertEmployeeOrStudent(db *mongo.Database, dataAbsen *models.Students_Employees) error {
	_, err := db.Collection("student_employees").InsertOne(context.TODO(), dataAbsen)
	if err != nil {
		log.Printf("Error Insert Data Absen: %v", err)
		return fmt.Errorf("Error inserting data: %v", err)
	}
	return nil
}
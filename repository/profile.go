package repository

import (
	"Ayala-Crea/server/config"
	"Ayala-Crea/server/models"
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertProfile(db *mongo.Database, profile *models.Profile) error {
	log.Println("Inserting profile into database...")
	_, err := db.Collection("profile").InsertOne(context.TODO(), profile)
	if err != nil {
		log.Printf("Error inserting profile: %v", err)
		return fmt.Errorf("could not insert profile: %v", err)
	}
	log.Println("Profile inserted successfully.")
	return nil
}

func UploadToFirebaseStorage(fileBytes []byte, filename string) (string, error) {
	if config.FirebaseStorageClient == nil {
		log.Println("Firebase Storage client is not initialized.")
		return "", fmt.Errorf("Firebase Storage client is not initialized")
	}

	ctx := context.Background()
	bucket := config.FirebaseStorageClient.Bucket(config.FirebaseBucketName)
	objectPath := "image-profile/" + filename
	object := bucket.Object(objectPath)
	writer := object.NewWriter(ctx)

	// Logging untuk debugging
	log.Println("Creating object in Firebase Storage...")
	log.Printf("File size: %d bytes", len(fileBytes))
	log.Printf("Filename: %s", filename)
	log.Printf("Object path: %s", objectPath)

	writer.ContentType = "image/jpeg" // Adjust the content type if needed
	writer.CacheControl = "public, max-age=86400"

	if _, err := writer.Write(fileBytes); err != nil {
		log.Printf("Error writing to the object: %v", err)
		return "", fmt.Errorf("cannot write to the object: %v", err)
	}

	if err := writer.Close(); err != nil {
		log.Printf("Error closing the writer: %v", err)
		return "", fmt.Errorf("cannot close the writer: %v", err)
	}

	log.Println("Object created successfully. Setting ACL...")

	// Mengatur ACL agar file dapat diakses publik
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		log.Printf("Error setting ACL: %v", err)
		return "", fmt.Errorf("cannot set ACL for the object: %v", err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", config.FirebaseBucketName, objectPath)
	log.Printf("File successfully uploaded to: %s", url)
	return url, nil
}

func GetProfileById(db *mongo.Database, id string) (*models.Profile, error) {
	var profile models.Profile

	filter := bson.M{"_id": id}
	err := db.Collection("profile").FindOne(context.TODO(), filter).Decode(&profile)
	if err != nil {
		return nil, fmt.Errorf("could not find prfile: %v", err)
	}
	return &profile, nil
}
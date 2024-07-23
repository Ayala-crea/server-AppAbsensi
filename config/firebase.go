package config

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var FirebaseStorageClient *storage.Client
var FirebaseBucketName = "appabsensi-6c600.appspot.com"

func InitFirebase() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("D:/semester 5/Proyek 3/server/appabsensi-6c600-firebase-adminsdk-bxg7i-9291f7cc77.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	FirebaseStorageClient = client
	log.Println("Firebase Storage client initialized successfully")
}

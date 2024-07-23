package repository

import (
	"Ayala-Crea/server/models"
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *mongo.Database, user *models.Users) error {
	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Jika id_role tidak diisi, atur nilainya ke 2 (atau nilai default yang diinginkan)
	if user.IdRole == 0 {
		user.IdRole = 1 // atau nilai default yang diinginkan
	}

	// Set IdUser jika belum diatur
	if user.IdUser.IsZero() {
		user.IdUser = primitive.NewObjectID()
	}

	// Simpan user ke database
	_, err = db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}
	return nil
}

func GetAdminByUsername(db *mongo.Database, username string) (*models.Users, error) {
	var user models.Users

	filter := bson.M{"username": username}
	err := db.Collection("users").FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}

	return &user, nil
}

func GetAdminByEmail(db *mongo.Database, email string) (*models.Users, error) {
	var user models.Users

	filter := bson.M{"email": email}
	err := db.Collection("users").FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}

	return &user, nil
}

func GetAdminById(db *mongo.Database, id primitive.ObjectID) (*models.Users, error) {
	var user models.Users

	filter := bson.M{"_id": id}
	err := db.Collection("users").FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}

	return &user, nil
}

func GetAdminDataById(db *mongo.Database, id primitive.ObjectID) (*models.Users, error) {
	var user models.Users

	filter := bson.M{"_id": id}
	err := db.Collection("users").FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}

	return &user, nil
}

func GenerateToken(user *models.Users) (string, error) {
	claims := &models.JWTClaims{
		IdUser: user.IdUser.Hex(), // Convert ObjectID to Hex string
		IdRole: user.IdRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(), // Token berlaku selama 1 jam
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
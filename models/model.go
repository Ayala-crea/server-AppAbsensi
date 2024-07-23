package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	IdUser      primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	IdRole      int                `bson:"id_role,omitempty" json:"id_role"`
	Nama        string             `bson:"nama,omitempty" json:"nama"`
	Username    string             `bson:"username,omitempty" json:"username"`
	Password    string             `bson:"password,omitempty" json:"password"`
	Email       string             `bson:"email,omitempty" json:"email"`
	PhoneNumber string             `bson:"phone_number,omitempty" json:"phone_number"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"-"`
	IsActive    string             `bson:"is_active,omitempty" json:"is_active"`
}

type Roles struct {
	IdRole int    `bson:"id_role,omitempty" json:"id_role"`
	Nama   string `bson:"nama,omitempty" json:"nama"`
}

type JWTClaims struct {
	jwt.StandardClaims
	IdUser string `json:"id_user"`
	IdRole int    `json:"id_role"`
}

type Profile struct {
	IdProfile primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	IdUser    string             `bson:"user_id,omitempty" json:"user_id"`
	Image     string             `bson:"img"`
	Address   string             `bson:"address,omitempty" json:"address"`
}

type Students_Employees struct {
	IdStudent_Employee primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	IdUser             uint               `bson:"id_user,omitempty" json:"id_user"`
	FullName           string             `bson:"full_name,omitempty" json:"full_name"`
	Status             string             `bson:"status,omitempty" json:"status"`
	Class              string             `bson:"class,omitempty" json:"class"`
	NpkOrNpm           string             `bson:"npk_or_npm,omitempty" json:"npk_or_npm"`
	PhoneNumber        string             `bson:"phone_number,omitempty" json:"phone_number"`
}

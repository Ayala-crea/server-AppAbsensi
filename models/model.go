package models

import (
    "time"

    "github.com/dgrijalva/jwt-go"
)

type Users struct {
    ID            int       `json:"id"`
    IdRole        int       `json:"id_role"`
    IdPenginputan string    `json:"id_penginputan"`
    Nama          string    `json:"nama"`
    Username      string    `json:"username"`
    Password      string    `json:"password"`
    Email         string    `json:"email"`
    PhoneNumber   string    `json:"phone_number"`
    CreatedAt     time.Time `json:"created_at"`
    IsActive      bool    `json:"is_active"`
}

type Roles struct {
    ID    int    `json:"id"`
    Nama  string `json:"nama"`
}

type JWTClaims struct {
    jwt.StandardClaims
    IdUser int `json:"id_user"`
    IdRole int `json:"id_role"`
}

type Profile struct {
    ID      int    `json:"id"`
    IdUser  int    `json:"user_id"`
    Image   string `json:"img"`
    Address string `json:"address"`
}

type StudentsEmployees struct {
    ID             int    `json:"id"`
    IdUser         int    `json:"id_user"`
    FullName       string `json:"full_name"`
    Email          string `json:"email"`
    Status         string `json:"status"`
    Class          string `json:"class"`
    NpkOrNpm       string `json:"npk_or_npm"`
    PhoneNumber    string `json:"phone_number"`
}

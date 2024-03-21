package models

import "github.com/dgrijalva/jwt-go"

type AdminModel struct {
	AdminID  int    `gorm:"primaryKey" `
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

package models

import "github.com/dgrijalva/jwt-go"

// user entity creation

type UsersModel struct {
	UserID   int    `gorm:"primaryKey;autoIncrement"`
	Phone    string `json:"phone" validate:"required"`
	Username string `json:"username"`
}

//verify the otp

type VerifyOTP struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Otp      string `json:"otp"`
}

//Userclaims struct for JWT authentication

type UserClaims struct {
	UserID uint 
	Phone string `json:"phone"`
	jwt.StandardClaims
}

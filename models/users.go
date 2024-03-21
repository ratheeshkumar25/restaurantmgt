package models

import "github.com/dgrijalva/jwt-go"

// user entity creation

type UsersModel struct {
	UserID   int    `gorm:"primaryKey" `
	Phone    string `json:"phone" validate:"required"`
	Username string `json:"name"`
	// 	Token         string `json:"token"`
	// 	Refresh_token string `json:"Refreshtoken"`
}

//verify the otp

type VerifyOTP struct {
	Phone string `json:"phone"`
	Otp   string `json:"otp"`
}

//Userclaims struct for JWT authentication

type UserClaims struct {
	jwt.StandardClaims
	Phone string `json:"phone"`
}

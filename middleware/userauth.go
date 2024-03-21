package middleware

import (
	"errors"
	"fmt"
	"restaurant/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin"
)

func Trim(token string)(string,error){
	parts := strings.SplitN(token,"Bearer ",2)
	return AuthenticateUser(parts[1])
}


func GenerateUsertoken(phone string) (string, error) {
	claims := jwt.MapClaims{
		"phone": phone,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil{
		return "",err
	}
	return tokenString,nil

}

func AuthenticateUser (signedStringToken string)(string,error){
	var userClaims models.UserClaims
	token, err := jwt.ParseWithClaims(signedStringToken, &userClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil // Replace with your secret key
	  })

	  if err != nil {
		return "",err
		}
	//check the token is valid 
	if !token.Valid{
		return "",errors.New("token is not valid")
	}
	//type assert the claims from the token object 
	claims,ok := token.Claims.(*models.UserClaims)

	if !ok {
		err = errors.New("couldn't parse claims")
		return "", err
	}
	phone := claims.Phone
	
	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("token expired")
		return "", err
	}

	return phone, nil
	
}	  

func UserauthMiddleware() gin.HandlerFunc {
		return func (c *gin.Context){
			// Extract token from the request header or other sources
			tokenString := c.GetHeader("Authorization")
			fmt.Println("Authorization Header",tokenString)
	
			// Check if token exists
			if tokenString == "" {
				c.AbortWithStatusJSON(401, gin.H{"error": "User Authorization is missing"})
				return
			}
		
			// Trim the token to get the actual token string
			authHeader := strings.TrimSpace(strings.TrimPrefix(tokenString,"Bearer"))
		
			phone, err := AuthenticateUser(authHeader)
			if err != nil {
				fmt.Println("Error authenticating user:", err)
				c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
				return
			}
		
			fmt.Println("Authenticated user:", phone)

		}
}






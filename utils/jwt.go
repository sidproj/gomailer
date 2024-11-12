package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "Hello world"

func GenerateJWT(userID string)(string,error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"user_id":userID,
		"exp":time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString,err := token.SignedString([]byte(secretKey))
	if err != nil{
		return "",err
	}
	return tokenString,nil
}

func VerifyJWT(JWTToken string) (string,error){
	token,err := jwt.Parse(JWTToken,func(token *jwt.Token)(interface {},error){
		return []byte(secretKey),nil
	})

	if err != nil{
		return "",err
	}

	if !token.Valid {
		return "",fmt.Errorf("invalid token")
	}
	return token.Claims.GetSubject()
}
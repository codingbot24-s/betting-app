package helpers

import (
	"log"
	"os"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)



func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

// verify the jwt then extract role from token
func VerifyToken(token string) (string, error) {
	var claims jwt.MapClaims
	fmt.Println("token parsing started")
	tokenClaims, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetEnv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	if !tokenClaims.Valid {
		return "", err
	}
	return claims["role"].(string), nil
}




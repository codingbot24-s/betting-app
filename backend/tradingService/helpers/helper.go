package helpers

import (
	"errors"

	"log"

	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}

func ValidateToken(tokenstr string) (*jwt.Token, error) {
	return jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetEnv("JWT_SECRET")), nil
	})
}

func GetUserIDFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}
	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid token")
	}

	return userId, nil

}

type OrderType string

const (
	Buy  OrderType = "buy"
	Sell OrderType = "sell"
)

type Order struct {
	Type OrderType
}

type SideType string

const (
	Yes SideType = "YES"
	No  SideType = "NO"
)

type Side struct {
	Type SideType
}

func CheckSideType(sideType SideType) bool {
	switch sideType {
	case Yes, No:
		return true
	default:
		return false
	}
}

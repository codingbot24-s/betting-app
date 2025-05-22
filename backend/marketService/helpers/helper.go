package helpers

import (
	"log"
	"os"
	"time"

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

const TimeFormat = "2006-01-02 15:04"

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(TimeFormat, timeStr)
}

func FormatTime(t time.Time) string {
	return t.Format(TimeFormat)
}

package helpers

import (
	"encoding/json"
	"errors"
	"time"

	"io"
	"log"
	"net/http"
	"net/url"
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


func ValidateToken(tokenstr string ) (*jwt.Token,error){
	return jwt.Parse(tokenstr,func(token * jwt.Token)(interface{},error){
		return []byte(GetEnv("JWT_SECRET")),nil
	})
}

func GetUserIDFromToken(token *jwt.Token) (string,error) {
	claims,ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}
	userId,ok := claims["user_id"].(string) 
	if !ok{
		return "", errors.New("invalid token")
	} 
		
	return userId,nil

}

func GetTeamSlug() []string {
	// return will be array of slugs
	// by making a rest call to the team service api/teams/slug
	baseurl := GetEnv("TEAM_SERVICE_URL_FOR_FETCHING_SLUG")
	// Parse the url
	parsedUrl, err := url.Parse(baseurl)

	if err != nil {
		log.Fatal("Error parsing URL:", err)
		return nil
	}

	resp, err := http.Get(parsedUrl.String())

	if err != nil {
		log.Fatal("Error making GET request:", err)
		return nil
	}
	defer resp.Body.Close()

	body,err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Error reading response body:", err)
		return nil
	}

		

	// Unmarshal the response body into a slice of strings
	var slugs []string
	err = json.Unmarshal(body, &slugs)	
	if err != nil {
		log.Fatal("Error unmarshalling response body:", err)
		return nil
	}
	return slugs
}

type OrderType string

const (
	Buy  OrderType = "buy"
	Sell OrderType = "sell"
)

type Order struct {
	Type OrderType 
}

// Check the Order Type

func CheckOrderType(orderType OrderType) bool {
	switch orderType {
	case Buy, Sell:
		return true
	default:
		return false
	}
}



func retry(attempts int, sleep time.Duration, fn func() error) error {
	err := fn()
	if err == nil {
		return nil
	}

	if attempts--; attempts > 0 {
		time.Sleep(sleep)
		return retry(attempts, sleep*2, fn)
	}

	return err
}
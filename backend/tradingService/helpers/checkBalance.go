package helpers

import (
	"encoding/json"
	
	"io"

	"net/http"
	"net/url"
)

// Pass The userID in query params to the user service by making a rest call and
// fetch the user balance
func CheckBalance(userID string) (float64, error) {
	
	// Get the balance from the user service
	// load url from .env file

	
	requrl := GetEnv("CHECK_BALANCE_URL")

	// Parse the url 
	parsedUrl, err := url.Parse(requrl)
	if err != nil {
		return 0, err
	}
	// add the query params
	query := parsedUrl.Query()
	query.Set("user_id", userID) 
	parsedUrl.RawQuery = query.Encode()

	// make the api call
	
	resp, err := http.Get(parsedUrl.String())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	
	var balance float64
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// Unmarshal the response body
	
	err = json.Unmarshal(body, &balance)
	if err != nil {
		
		return 0, err
	}

	return balance, nil
}



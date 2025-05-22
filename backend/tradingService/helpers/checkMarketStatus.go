package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	
)

type MarketResponse struct {
	IsOpen bool `json:"isOpen"`
}

func CheckMarketStatus(marketID string) (bool, error) {

	resp, err := http.Get(fmt.Sprintf("http://localhost:8083/market/%s", marketID))
	if err != nil {
		return false, fmt.Errorf("failed to check market status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("market service returned status: %d", resp.StatusCode)
	}


	var market MarketResponse
	if err := json.NewDecoder(resp.Body).Decode(&market); err != nil {
		return false, fmt.Errorf("failed to decode market response: %v", err)
	}


	// market resonse is a boolean
	if market.IsOpen {
		return true, nil
	} else {
		return false, nil
	}
}

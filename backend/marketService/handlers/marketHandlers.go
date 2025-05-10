package handlers

import (
	"fmt"
	"net/http"
)

// TODO: implement the admin checker middleware

type CreateMarketReq struct {

	Question string `json:"question"`
	Category string `json:"category"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	Status string `json:"status"`
	Outcome *bool `json:"outcome"`
}

func CreatedMarketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside CreatedMarketHandler")
}
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"slices"

	"github.com/codingbot24-s/common"
	"github.com/codingbot24-s/db"
	"github.com/codingbot24-s/helpers"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/datatypes"
)

// create a outcome handeler for market also

type CreateMarketReq struct {
	Question  string `json:"question"`
	Category  string `json:"category"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
	Outcome   *bool  `json:"outcome"`
}

type MarketResponse struct {
	ID        string `json:"id"`
	Question  string `json:"question"`
	Category  string `json:"category"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
	Outcome   *bool  `json:"outcome"`
}

func marketToResponse(m db.Market) MarketResponse {
	return MarketResponse{
		ID:        m.ID,
		Question:  m.Question,
		Category:  m.Category,
		StartTime: helpers.FormatTime(m.StartTime),
		EndTime:   helpers.FormatTime(m.EndTime),
		Status:    string(m.Status),
		Outcome:   m.Outcome,
	}
}

var validate *validator.Validate

func CreatedMarketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Parse the request body
	var reqBody CreateMarketReq
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request body
	validate = validator.New()
	err = validate.Struct(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Check if the start time is before the end time
	if reqBody.StartTime >= reqBody.EndTime {
		http.Error(w, "Start time must be before end time", http.StatusBadRequest)
		return
	}
	// Check if the status is valid
	validStatuses := []string{"draft", "open", "closed", "resolved"}
	isValidStatus := slices.Contains(validStatuses, reqBody.Status)

	if !isValidStatus {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	// parse the time

	startTime, err := helpers.ParseTime(reqBody.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}
	endTime, err := helpers.ParseTime(reqBody.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	market := db.Market{
		Question:  reqBody.Question,
		Category:  reqBody.Category,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    common.MarketStatus(reqBody.Status),
		Outcome:   reqBody.Outcome,
	}

	// Save the market to the database
	err = db.DB.Create(&market).Error
	if err != nil {
		http.Error(w, "Failed to create market", http.StatusInternalServerError)
		return
	}

	// Return the created market as a response
	w.WriteHeader(http.StatusCreated)

	reponse := marketToResponse(market)

	json.NewEncoder(w).Encode(reponse)

}

func ListActiveMarketsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	markets := []db.Market{}
	err := db.DB.Where("status = ?", common.StatusOpen).Find(&markets).Error
	if err != nil {
		http.Error(w, "Failed to list markets", http.StatusInternalServerError)
		return
	}

	marketsResponse := make([]MarketResponse, len(markets))
	for i, market := range markets {
		marketsResponse[i] = marketToResponse(market)
	}

	response := map[string]interface{}{
		"message": "Markets listed successfully",
		"markets": marketsResponse,
	}

	json.NewEncoder(w).Encode(response)
}

func ListClosedMarketsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	markets := []db.Market{}
	err := db.DB.Where("status = ?", common.StatusClosed).Find(&markets).Error
	if err != nil {
		http.Error(w, "Failed to list markets", http.StatusInternalServerError)
		return
	}

	marketsResponse := make([]MarketResponse, len(markets))
	for i, market := range markets {
		marketsResponse[i] = marketToResponse(market)
	}

	response := map[string]interface{}{
		"message": "Closed markets listed successfully",
		"markets": marketsResponse,
	}

	json.NewEncoder(w).Encode(response)

}

type ResolvedMarketReq struct {
	Outcome *bool `json:"outcome"`
}


func ResolvedMarketsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	marketid := mux.Vars(r)["id"]

	// Parse the request body
	var reqBody ResolvedMarketReq
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the market from the database
	var market db.Market
	err = db.DB.Where("id = ?", marketid).First(&market).Error
	if err != nil {
		http.Error(w, "Failed to get market", http.StatusInternalServerError)
		return
	}

	if market.Status != common.StatusClosed {
		http.Error(w, "Market is not closed So cannot be resolved", http.StatusBadRequest)
		return
	}

	if market.Outcome != nil {
		http.Error(w, "Market is already resolved so cant be resolved again", http.StatusBadRequest)
		return
	}

	market.Outcome = reqBody.Outcome
	// agr status resolved kardiya to open or closed route pr request se vo market nahi milega because market status is resolved so naya handler banana padega only for resolved market
	market.Status = common.StatusResolved
	// begin transaction from here
	tx := db.DB.Begin()
	err = tx.Save(market).Error
	if err != nil {
		http.Error(w, "Failed to save market", http.StatusInternalServerError)
		tx.Rollback()

	}
	payload := map[string]interface{}{
		"marketId": market.ID,
		"outcome":  market.Outcome,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
		tx.Rollback()
	}

	event := db.OutboxEvent{
		ID:        uuid.New(),
		EventType:  "market-resolved",
		Payload:   datatypes.JSON(data),
		Processed: false,
		CreatedAt: time.Now(),
	}

	err = tx.Create(&event).Error
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		tx.Rollback()
	}

	err = tx.Commit().Error
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
	}
}

func SendMarketStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		marketid := mux.Vars(r)["id"]

		var market db.Market
		err := db.DB.Where("id = ?", marketid).First(&market).Error
		if err != nil {
			http.Error(w, "Failed to get market", http.StatusInternalServerError)
			return
		}

		marketResponse := marketToResponse(market)

		response := map[string]bool{
			"isOpen": marketResponse.Status == "open",
		}

		json.NewEncoder(w).Encode(response)
	}
}

func CloseMarketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	marketid := mux.Vars(r)["id"]
	var market db.Market
	err := db.DB.Where("id = ?", marketid).First(&market).Error
	if err != nil {
		http.Error(w, "Failed to get market", http.StatusInternalServerError)
		return
	}		
	if market.Status != common.StatusOpen {
		http.Error(w, "Market is not open So cannot be closed", http.StatusBadRequest)
		return
	}
	market.Status = common.StatusClosed
	err = db.DB.Save(&market).Error
	if err != nil {
		http.Error(w, "Failed to save market", http.StatusInternalServerError)
		return
	}
	
}
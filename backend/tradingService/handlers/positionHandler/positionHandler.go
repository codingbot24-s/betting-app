package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/codingbot24-s/db"
	"github.com/codingbot24-s/db/models"
	"github.com/codingbot24-s/helpers"
	"github.com/codingbot24-s/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	SideYes = "YES"
	SideNo  = "NO"
)

type Position struct {
	UserID   string  `json:"userId" gorm:"type:uuid;not null"`
	MarketID string  `json:"marketId" gorm:"type:uuid;not null"`
	Side     string  `json:"side" gorm:"not null" validate:"required,oneof=YES NO"`
	Amount   float64 `json:"amount" gorm:"not null" validate:"required,gt=0"`
}

type PositionResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func CreatePositionHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Request recieived to place order")
	var reqBody Position
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the order
	validate := validator.New()
	err = validate.Struct(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate the side

	if !helpers.CheckSideType(helpers.SideType(reqBody.Side)) {
		http.Error(w, "Invalid side", http.StatusBadRequest)
		return
	}
	if reqBody.Side != SideYes && reqBody.Side != SideNo {
		http.Error(w, "Side must be either YES or NO", http.StatusBadRequest)
		return
	}

	// Check the validity of amount
	if reqBody.Amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	reqBody.UserID = userID

	marketOpen, err := helpers.CheckMarketStatus(reqBody.MarketID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !marketOpen {
		http.Error(w, "Market is not open", http.StatusBadRequest)
		return
	}
	// TODO: check the balance later
	// fmt.Println("checking is user has enough funds")
	// balance, err := helpers.CheckBalance(userID)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// if balance < reqBody.Amount {
	// 	http.Error(w, "Insufficient balance", http.StatusBadRequest)
	// 	return
	// }
	// create the position for the user
	position := models.Position{
		ID:        uuid.New(),
		UserID:    userID,
		MarketID:  reqBody.MarketID,
		Side:      reqBody.Side,
		Amount:    reqBody.Amount,
		CreatedAt: time.Now(),
	}
	p, err := helpers.CreatePosition(db.DB, position)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)

	fmt.Println("Order placed successfully")
}

func GetUserPositionsHandler(w http.ResponseWriter, r *http.Request) {

	userid := mux.Vars(r)["userid"]

	positions, err := helpers.GetUserPositions(db.DB, userid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(positions)

}

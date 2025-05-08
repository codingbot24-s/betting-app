package tradingHandlers

import (
	"encoding/json"
	"fmt"

	"net/http"

	"slices"

	"github.com/codingbot24-s/db"
	"github.com/codingbot24-s/db/models"
	"github.com/codingbot24-s/helpers"
	"github.com/codingbot24-s/middlewares"
	"github.com/go-playground/validator/v10"
)

// OrderType is a custom type for order types
// TODO: Add a sell function also
func Buy(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(string)

	w.Header().Set("Content-Type", "application/json")
	// parse the request body
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate the order
	validate := validator.New()
	err = validate.Struct(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate the order type
	if !helpers.CheckOrderType(helpers.OrderType(order.OrderType)) {
		http.Error(w, "Invalid order type", http.StatusBadRequest)
		return
	}

	// check the quantity
	if order.Quantity <= 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	order.UserID = userID

	// Calculate total amount before balance check
	order.Total = float64(order.Quantity) * order.Price

	// check if the user has enough balance make a rest call to the user service
	balance, err := helpers.CheckBalance(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if balance < order.Total {
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}
	// get the team slugs from team service
	teamSlug := helpers.GetTeamSlug()
	isValidTeam := slices.Contains(teamSlug, order.Team)

	if !isValidTeam {
		http.Error(w, "Invalid team", http.StatusBadRequest)
		return
	}

	// now place the order with transaction
	fmt.Println("Placing order")
	err = helpers.CreateOrder(db.DB, order)
	if err != nil {
		http.Error(w, "Creating order failed", http.StatusInternalServerError)
		return
	}
	
	fmt.Println("Order placed successfully")
	
	//TODO: after order is placed update the user holdings also 
}


// TODO: Holding check for user to sell the stock 
// func Sell(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	// parse the request body

// 	var order models.Order
// 	err := json.NewDecoder(r.Body).Decode(&order)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// validate the order
// 	validate := validator.New()
// 	err = validate.Struct(&order)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	// validate the order type
// 	if !helpers.CheckOrderType(helpers.OrderType(order.OrderType)) {
// 		http.Error(w, "Invalid order type", http.StatusBadRequest)
// 		return
// 	}

// 	// check the quantity
// 	if order.Quantity <= 0 {
// 		http.Error(w, "Invalid quantity", http.StatusBadRequest)
// 		return
// 	}

// 	// TODO: Check if the user has enough stocks to sell
// 	userID := r.Context().Value(middlewares.UserIDKey).(string)


// 	stocks, err := helpers.CheckStocks(userID, order.Team, order.Quantity)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if stocks < order.Quantity {
// 		http.Error(w, "Insufficient stocks", http.StatusBadRequest)
// 		return
// 	}
// }
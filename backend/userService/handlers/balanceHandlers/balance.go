package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"time"

	"github.com/codingbot24-s/db"
	modles "github.com/codingbot24-s/db/models"
	"github.com/codingbot24-s/middlewares"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func GetBalancePrivate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)

	var user modles.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user.Balance)
}

type AddBalanceRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0,lte=1000000"`
}

// when adding balance, create a transaction and store it in the database if the transaction is successful then return the balance if not then return the error and not increase the balance

func AddBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)

	var user modles.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var addBalanceRequest AddBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&addBalanceRequest); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validate.Struct(addBalanceRequest); err != nil {
		http.Error(w, "Invalid amount: must be greater than 0 and less than 1,000,000", http.StatusBadRequest)
		return
	}

	// Start a database transaction
	tx := db.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Update user balance
	user.Balance += addBalanceRequest.Amount
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update balance", http.StatusInternalServerError)
		return
	}

	// Create transaction record
	transaction := modles.Transaction{
		User:            user,
		UserID:          user.ID,
		Amount:          addBalanceRequest.Amount,
		Description:     "Balance deposit",
		TransactionDate: time.Now(),
		TransactionID:   fmt.Sprintf("DEP-%d-%d", user.ID, time.Now().UnixNano()),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to record transaction", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to complete transaction", http.StatusInternalServerError)
		return
	}

	// Return updated balance
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance":        user.Balance,
		"transaction_id": transaction.TransactionID,
	})
}

type WithdrawBalanceRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0,lte=1000000"`
}


// Call this function to withdraw from the user's account
func WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)
  // find the user by id 
	var user modles.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
  // will take an Amount to withdraw 
	var withdrawBalanceRequest WithdrawBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&withdrawBalanceRequest); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validate.Struct(withdrawBalanceRequest); err != nil {
		http.Error(w, "Invalid amount: must be greater than 0 and less than 1,000,000", http.StatusBadRequest)
		return
	}

	// Check if the balance is enough
	if user.Balance < withdrawBalanceRequest.Amount {
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	// Start a database transaction
	tx := db.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Update user balance
	user.Balance -= withdrawBalanceRequest.Amount
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update balance", http.StatusInternalServerError)
		return
	}

	// Create transaction record
	transaction := modles.Transaction{
		User:            user,
		UserID:          user.ID,
		Amount:          -withdrawBalanceRequest.Amount, // Negative amount for withdrawal
		Description:     "Balance withdrawal",
		TransactionDate: time.Now(),
		TransactionID:   fmt.Sprintf("WD-%d-%d", user.ID, time.Now().UnixNano()),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to record transaction", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to complete transaction", http.StatusInternalServerError)
		return
	}

	// Return updated balance and transaction info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance":        user.Balance,
		"transaction_id": transaction.TransactionID,
	})
}


type TransactionResponse struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	UserID          uint      `json:"user_id"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
	TransactionID   string    `json:"transaction_id"`
}

type TransactionHistoryResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                        `json:"total"`
}

func GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)

	var user modles.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Get total count
	var total int64
	if err := db.DB.Model(&modles.Transaction{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		http.Error(w, "Failed to get transaction count", http.StatusInternalServerError)
		return
	}

	// Get all transactions
	var transactions []modles.Transaction
	if err := db.DB.Where("user_id = ?", userID).Order("transaction_date desc").Find(&transactions).Error; err != nil {
		http.Error(w, "Failed to get transaction history", http.StatusInternalServerError)
		return
	}

	
	transactionResponses := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		transactionResponses[i] = TransactionResponse{
			ID:              t.ID,
			CreatedAt:       t.CreatedAt,
			UpdatedAt:       t.UpdatedAt,
			UserID:          t.User.ID,
			Amount:          t.Amount,
			Description:     t.Description,
			TransactionDate: t.TransactionDate,
			TransactionID:   t.TransactionID,
		}
	}

	// Prepare response
	response := TransactionHistoryResponse{
		Transactions: transactionResponses,
		Total:        total,
	}

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func GetBalancePublic(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		http.Error(w, "User ID is required in the query params", http.StatusBadRequest)
		return
	}

	// 
	var user modles.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	
	json.NewEncoder(w).Encode(user.Balance)
}

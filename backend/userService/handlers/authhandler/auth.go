package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/codingbot24-s/db"
	modles "github.com/codingbot24-s/db/models"
	"github.com/codingbot24-s/helpers"
	"github.com/codingbot24-s/middlewares"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type RegisterResponse struct {
	Token string      `json:"token"`
	User  modles.User `json:"user"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user modles.User

	// Add request body logging
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	fmt.Printf("Raw request body: %s\n", string(body))

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Decode error: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Printf("Decoded user: %+v\n", user)

	if err := validate.Struct(&user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed:\n"
		for _, e := range validationErrors {
			errorMsg += fmt.Sprintf("Field: %s, Tag: %s, Value: %v\n",
				e.Field(), e.Tag(), e.Value())
		}
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	fmt.Println("user after validation", user)

	var existingUser modles.User
	err = db.DB.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser).Error
	if err == nil {
		http.Error(w, "Username or email already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	err = db.DB.Create(&user).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := helpers.GenerateToken(strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{
		Token: token,
		User:  user,
	}

	json.NewEncoder(w).Encode(response)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  modles.User `json:"user"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// first validate the token then extract userid from it then check if user exists in db then compare the password if all ok then return the response

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(loginReq); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	// get the token from auth header
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	// validate the token
	jwtToken, err := helpers.ValidateToken(authToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// extract userid from token
	userid, err := helpers.GetUserIDFromToken(jwtToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	var user modles.User
	if err := db.DB.Where("id = ?", userid).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := helpers.ComparePassword(user.Password, loginReq.Password); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	response := LoginResponse{
		Token: authToken,
		User:  user,
	}

	json.NewEncoder(w).Encode(response)
}

func GetSingleUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// middleware will pass the userid to the request
	userID := r.Context().Value(middlewares.UserIDKey).(string)

	var user modles.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}

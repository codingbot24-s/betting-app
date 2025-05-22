package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/codingbot24-s/db"

	"github.com/go-playground/validator/v10"
)

func GetTeams(w http.ResponseWriter, r *http.Request) {

	var teams []db.Team
	err := db.DB.Find(&teams).Error
	if err != nil {
		http.Error(w, "Error fetching teams", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(teams)
}

// temporary function to create a Teams in db it will be replaced by the live match service

func CreateTeam(w http.ResponseWriter, r *http.Request) {
	var teams []db.Team
	err := json.NewDecoder(r.Body).Decode(&teams)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// verify each team is valid
	validate := validator.New()
	for _, team := range teams {
		err = validate.Struct(team)
		if err != nil {
			http.Error(w, "Error validating team: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Create all teams
	for _, team := range teams {
		if db.DB.Create(&team).Error != nil {
			http.Error(w, "Error creating team", http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(teams)
}

func GetTeamSlug(w http.ResponseWriter, r *http.Request) {
	// fetch the all teams symbols from db
	var teams []db.Team
	err := db.DB.Find(&teams).Error
	if err != nil {
		http.Error(w, "Error fetching teams", http.StatusInternalServerError)
		return
	}

	// create an array of team symbols
	teamSymbols := []string{}
	for _, team := range teams {
		teamSymbols = append(teamSymbols, team.TeamSymbol)
	}
	json.NewEncoder(w).Encode(teamSymbols)
}

package routes
	
import (
	
	"github.com/gorilla/mux"
	"github.com/codingbot24-s/handlers"
)

func SetupTeamRoutes(router *mux.Router) *mux.Router {
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/teams", handlers.CreateTeam).Methods("POST")
	api.HandleFunc("/teams", handlers.GetTeams).Methods("GET")
	api.HandleFunc("/teams/slug", handlers.GetTeamSlug).Methods("GET")
	return router
}




package server

import (
	"github.com/gorilla/handlers"
	"net/http"
)

func GetCORSConfig() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                      
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), 
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),          
		handlers.AllowCredentials(),                                                  
	)
}

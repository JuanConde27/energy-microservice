package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/JuanConde27/energy-microservice/src/config"
)

func Start() {
	config.Migrate()

	router := mux.NewRouter()

	GetCORSConfig()

	server := http.Server{
		Addr:    ":3000",
		Handler: GetCORSConfig()(router),
	}

	log.Println(`Servidor escuchando en el puerto 3000`)
	log.Fatal(server.ListenAndServe())
}

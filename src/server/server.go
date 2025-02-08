package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/JuanConde27/energy-microservice/src/config"
	"github.com/JuanConde27/energy-microservice/src/utils"
)

func Start() {
	config.Migrate()

	csvPath := "test_bia.csv"
	utils.LoadCSVData(csvPath)

	router := mux.NewRouter()

	RegisterRoutes(router)

	GetCORSConfig()

	server := http.Server{
		Addr:    ":3000",
		Handler: GetCORSConfig()(router),
	}

	log.Println(`Servidor escuchando en el puerto 3000`)
	log.Fatal(server.ListenAndServe())
}

package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/JuanConde27/energy-microservice/src/config"
	"github.com/JuanConde27/energy-microservice/src/utils"
)

func SetupRouter() *mux.Router {
	config.Migrate()

	if os.Getenv("TEST_MODE") != "true" {
		csvPath := "/app/test_bia.csv"
		utils.LoadCSVData(csvPath)
	}

	router := mux.NewRouter()
	RegisterRoutes(router)

	return router
}

func Start() {
	router := SetupRouter()

	server := http.Server{
		Addr:    ":3000",
		Handler: GetCORSConfig()(router),
	}

	log.Println("Servidor escuchando en el puerto 3000")
	log.Fatal(server.ListenAndServe())
}

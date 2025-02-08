package routes

import (
	"github.com/JuanConde27/energy-microservice/src/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterConsumptionRoutes(router *mux.Router, controller *controllers.ConsumptionController) {
	router.HandleFunc("/consumption", controller.GetConsumption).Methods(http.MethodGet)
}

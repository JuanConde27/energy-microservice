package routes

import (
	"github.com/JuanConde27/energy-microservice/src/controllers"
	"github.com/JuanConde27/energy-microservice/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterConsumptionRoutes(router *mux.Router, controller *controllers.ConsumptionController) {
	router.Use(middlewares.RecoveryMiddleware)
	router.HandleFunc("/consumption", controller.GetConsumption).Methods(http.MethodGet)
}

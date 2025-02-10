package server

import (
	"github.com/gorilla/mux"
	"github.com/JuanConde27/energy-microservice/src/controllers"
	"github.com/JuanConde27/energy-microservice/src/repositories"
	"github.com/JuanConde27/energy-microservice/src/routes"
	"github.com/JuanConde27/energy-microservice/src/services"
	"github.com/JuanConde27/energy-microservice/src/config"
)

func RegisterRoutes(router *mux.Router) {
	db := config.GetConnection()

	consumptionRepo := repositories.NewConsumptionRepository(db)
	consumptionService := services.NewConsumptionService(consumptionRepo)
	consumptionController := controllers.NewConsumptionController(consumptionService)

	routes.RegisterConsumptionRoutes(router, consumptionController)
}

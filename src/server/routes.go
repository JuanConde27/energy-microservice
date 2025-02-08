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
	// Obtener la conexión a la base de datos
	db := config.GetConnection()

	// Inicializar repositorio, servicio y controlador
	consumptionRepo := repositories.NewConsumptionRepository(db)
	consumptionService := services.NewConsumptionService(consumptionRepo)
	consumptionController := controllers.NewConsumptionController(consumptionService)

	// Pasar el router y el controlador a la función de rutas
	routes.RegisterConsumptionRoutes(router, consumptionController)
}

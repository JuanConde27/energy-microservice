package controllers

import (
	"net/http"
	
	"github.com/JuanConde27/energy-microservice/src/config"
	"github.com/JuanConde27/energy-microservice/src/utils"
	"github.com/JuanConde27/energy-microservice/src/services"
)

type ConsumptionController struct {
	Service *services.ConsumptionService
}

func NewConsumptionController(service *services.ConsumptionService) *ConsumptionController {
	return &ConsumptionController{Service: service}
}

func (cc *ConsumptionController) GetConsumption(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()

    meterIDs, startDate, endDate, kindPeriod, statusCode, errorMessage := utils.ValidateQueryParams(query)
    if statusCode != http.StatusOK {
        http.Error(w, errorMessage, statusCode)
        return
    }

    response, err := cc.Service.GetConsumption(meterIDs, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), kindPeriod)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    config.SendResponse(w, http.StatusOK, response)
}
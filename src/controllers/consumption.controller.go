package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/JuanConde27/energy-microservice/src/services"
)

type ConsumptionController struct {
	Service *services.ConsumptionService
}

func NewConsumptionController(service *services.ConsumptionService) *ConsumptionController {
	return &ConsumptionController{Service: service}
}

func parseMeterIDs(param string) []int {
	ids := strings.Split(param, ",")
	var result []int
	for _, id := range ids {
		parsedID, _ := strconv.Atoi(id)
		result = append(result, parsedID)
	}
	return result
}

func (cc *ConsumptionController) GetConsumption(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	meterIDs := parseMeterIDs(query.Get("meters_ids"))
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")
	kindPeriod := query.Get("kind_period")

	response, err := cc.Service.GetConsumption(meterIDs, startDate, endDate, kindPeriod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/JuanConde27/energy-microservice/src/services"
)

type ConsumptionController struct {
	Service *services.ConsumptionService
}

func NewConsumptionController(service *services.ConsumptionService) *ConsumptionController {
	return &ConsumptionController{Service: service}
}

func parseMeterIDs(param string) ([]int, error) {
	ids := strings.Split(param, ",")
	var result []int
	for _, id := range ids {
		parsedID, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedID)
	}
	return result, nil
}

func isValidPeriod(kindPeriod string) bool {
	validPeriods := map[string]bool{"daily": true, "weekly": true, "monthly": true}
	return validPeriods[kindPeriod]
}

func isValidDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}

func (cc *ConsumptionController) GetConsumption(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	meterIDs, err := parseMeterIDs(query.Get("meters_ids"))
	if err != nil {
		http.Error(w, "❌ Invalid meter_ids format. Expected comma-separated integers.", http.StatusBadRequest)
		return
	}

	startDate, err := isValidDate(query.Get("start_date"))
	if err != nil {
		http.Error(w, "❌ Invalid start_date format. Expected YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	endDate, err := isValidDate(query.Get("end_date"))
	if err != nil {
		http.Error(w, "❌ Invalid end_date format. Expected YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	if startDate.After(endDate) {
		http.Error(w, "❌ start_date cannot be later than end_date.", http.StatusBadRequest)
		return
	}

	kindPeriod := query.Get("kind_period")
	if !isValidPeriod(kindPeriod) {
		http.Error(w, "❌ Invalid kind_period. Only 'daily', 'weekly', or 'monthly' are allowed.", http.StatusBadRequest)
		return
	}

	response, err := cc.Service.GetConsumption(meterIDs, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), kindPeriod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

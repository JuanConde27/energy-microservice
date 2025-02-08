package services

import (
	"github.com/JuanConde27/energy-microservice/src/models"
	"github.com/JuanConde27/energy-microservice/src/repositories"
	"time"
)

type ConsumptionService struct {
	Repo *repositories.ConsumptionRepository
}

func NewConsumptionService(repo *repositories.ConsumptionRepository) *ConsumptionService {
	return &ConsumptionService{Repo: repo}
}

func formatPeriod(t time.Time, period string) string {
	switch period {
	case "monthly":
		return t.Format("Jan 2006")
	case "weekly":
		return t.Format("Jan 2 - Jan 8")
	default:
		return t.Format("Jan 2")
	}
}

func (s *ConsumptionService) GetConsumption(meterIDs []int, startDate, endDate string, period string) (map[string]interface{}, error) {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var consumptions []models.Consumption  
	var err error

	switch period {
	case "monthly":
		consumptions, err = s.Repo.GetMonthlyConsumption(meterIDs, start, end)
	case "weekly":
		consumptions, err = s.Repo.GetWeeklyConsumption(meterIDs, start, end)
	case "daily":
		consumptions, err = s.Repo.GetDailyConsumption(meterIDs, start, end)
	}
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"period":    []string{},
		"data_graph": []map[string]interface{}{},
	}

	for _, c := range consumptions {
		response["period"] = append(response["period"].([]string), formatPeriod(c.Timestamp, period))
		data := map[string]interface{}{
			"meter_id":            c.MeterID,
			"address":             "Dirección mock",
			"active":              []float64{c.Consumption},
			"reactive_inductive":  []float64{0},
			"reactive_capacitive": []float64{0},
			"exported":            []float64{0},
		}
		response["data_graph"] = append(response["data_graph"].([]map[string]interface{}), data)
	}

	return response, nil
}

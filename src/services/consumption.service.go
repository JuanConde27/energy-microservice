package services

import (
    "strings"
    "time"
	"fmt"
	"sort"

    "github.com/JuanConde27/energy-microservice/src/repositories"
)

type ConsumptionResponse struct {
    Period    []string                 `json:"period"`
    DataGraph []ConsumptionDataGraph   `json:"data_graph"`
}

type ConsumptionDataGraph struct {
    MeterID            int       `json:"meter_id"`
    Address            string    `json:"address"`
    Active             []float64 `json:"active"`
    ReactiveInductive  []float64 `json:"reactive_inductive"`
    ReactiveCapacitive []float64 `json:"reactive_capacitive"`
    Exported           []float64 `json:"exported"`
}

type ConsumptionService struct {
    Repo *repositories.ConsumptionRepository
}

func NewConsumptionService(repo *repositories.ConsumptionRepository) *ConsumptionService {
    return &ConsumptionService{Repo: repo}
}

func formatPeriod(t time.Time, period string, startDate, endDate time.Time) string {
    if t.Before(startDate) || t.After(endDate) {
        return ""
    }
    switch period {
    case "monthly":
        return strings.ToUpper(t.Format("Jan 2006"))
    case "weekly":
		diffDays := int(t.Sub(startDate).Hours() / 24)
		weekIndex := diffDays / 7
		startOfWeek := startDate.AddDate(0, 0, weekIndex*7)
		endOfWeek := startOfWeek.AddDate(0, 0, 6)

		if endOfWeek.After(endDate) {
			endOfWeek = endDate
		}
	
		return strings.ToUpper(startOfWeek.Format("Jan 2")) + " - " + strings.ToUpper(endOfWeek.Format("Jan 2"))	
    case "daily":
        return strings.ToUpper(t.Format("Jan 2"))
    default:
        return strings.ToUpper(t.Format("Jan 2"))
    }
}

func nextPeriod(current time.Time, period string) time.Time {
    switch period {
    case "monthly":
        return current.AddDate(0, 1, 0)
    case "weekly":
        return current.AddDate(0, 0, 7)
    case "daily":
        return current.AddDate(0, 0, 1)
    default:
        return current.AddDate(0, 0, 1)
    }
}

func (s *ConsumptionService) GetConsumption(meterIDs []int, startDate, endDate string, period string) (ConsumptionResponse, error) {
    start, _ := time.Parse("2006-01-02", startDate)
    end, _ := time.Parse("2006-01-02", endDate)
    
    // Agregamos 23 horas, 59 minutos y 59 segundos para incluir todo el día final.
    end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

    if period == "weekly" {
        daysSinceStart := int(end.Sub(start).Hours() / 24)
        extraDays := (7 - (daysSinceStart % 7)) % 7
        end = end.AddDate(0, 0, extraDays)
    }

    consumptions, err := s.Repo.GetConsumptionByPeriod(meterIDs, start, end, period)
    if err != nil {
        return ConsumptionResponse{}, err
    }

    response := ConsumptionResponse{
        Period:    []string{},
        DataGraph: []ConsumptionDataGraph{},
    }

    dateMap := make(map[int]map[string]float64)
    for d := start; !d.After(end); d = nextPeriod(d, period) {
        dateStr := formatPeriod(d, period, start, end)

        if period == "weekly" {
            endOfWeek := d.AddDate(0, 0, 6)
            if endOfWeek.After(end) {
                continue 
            }
        }

        for _, meterID := range meterIDs {
            if _, exists := dateMap[meterID]; !exists {
                dateMap[meterID] = make(map[string]float64)
            }
            dateMap[meterID][dateStr] = 0
        }
        response.Period = append(response.Period, dateStr)
    }

    for _, c := range consumptions {
        periodStr := formatPeriod(c.Period, period, start, end)
        if periodStr == "" {
            continue
        }
    
        fmt.Println("🧐 Verificando formato: ", c.Period, "➡", periodStr)
    
        dateMap[c.MeterID][periodStr] = c.Consumption
    }
    
    // 🔍 Verificar valores en dateMap antes de construir la respuesta final
    fmt.Println("✅ Verificación final de dateMap:")
    for meterID, periods := range dateMap {
        for period, consumption := range periods {
            fmt.Println("Meter ID:", meterID, "Periodo:", period, "Consumo:", consumption)
        }
    }

    // Ordenamos las fechas en el periodo antes de llenar los datos en response
    sort.SliceStable(response.Period, func(i, j int) bool {
        dateI, _ := time.Parse("Jan 2", response.Period[i])
        dateJ, _ := time.Parse("Jan 2", response.Period[j])
        return dateI.Before(dateJ)
    })

    for _, meterID := range meterIDs {
        data := &ConsumptionDataGraph{
            MeterID:            meterID,
            Address:            "Dirección mock",
            Active:             []float64{},
            ReactiveInductive:  []float64{},
            ReactiveCapacitive: []float64{},
            Exported:           []float64{},
        }

        for _, date := range response.Period {
            fmt.Printf("📌 Asignando en response: Meter ID: %d Periodo: %s Consumo: %f\n", meterID, date, dateMap[meterID][date])
            data.Active = append(data.Active, dateMap[meterID][date])
            data.ReactiveInductive = append(data.ReactiveInductive, 0)
            data.ReactiveCapacitive = append(data.ReactiveCapacitive, 0)
            data.Exported = append(data.Exported, 0)
        }

        response.DataGraph = append(response.DataGraph, *data)
    }

    return response, nil
}

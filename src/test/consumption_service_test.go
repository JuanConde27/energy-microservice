package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/JuanConde27/energy-microservice/src/models"
	"github.com/JuanConde27/energy-microservice/src/services"
)

type FakeConsumptionRepository struct{}

func (f *FakeConsumptionRepository) GetConsumptionByPeriod(meterIDs []int, startDate, endDate time.Time, periodType string) ([]models.ConsumptionAggregate, error) {
	if periodType == "daily" {
		return []models.ConsumptionAggregate{
			{
				MeterID:     1,
				Consumption: 100.0,
				Period:      time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				MeterID:     1,
				Consumption: 200.0,
				Period:      time.Date(2023, 6, 2, 0, 0, 0, 0, time.UTC),
			},
			{
				MeterID:     1,
				Consumption: 300.0,
				Period:      time.Date(2023, 6, 3, 0, 0, 0, 0, time.UTC),
			},
		}, nil
	}
	return []models.ConsumptionAggregate{}, nil
}

func TestGetConsumptionServiceDaily(t *testing.T) {
	
	fakeRepo := &FakeConsumptionRepository{}
	service := services.NewConsumptionService(fakeRepo)


	meterIDs := []int{1}
	startDate := "2023-06-01"
	endDate := "2023-06-03"
	period := "daily"

	response, err := service.GetConsumption(meterIDs, startDate, endDate, period)
	assert.NoError(t, err)

	expectedPeriods := []string{"JUN 1", "JUN 2", "JUN 3"}
	assert.Equal(t, expectedPeriods, response.Period, "El arreglo de periodos debe coincidir con lo esperado")

	assert.Len(t, response.DataGraph, 1)
	dataGraph := response.DataGraph[0]
	assert.Equal(t, 1, dataGraph.MeterID, "El meter_id debe ser 1")
	assert.Equal(t, "Dirección mock", dataGraph.Address, "La dirección debe ser la de mock")

	expectedActive := []float64{100.0, 200.0, 300.0}
	assert.Equal(t, expectedActive, dataGraph.Active, "Los consumos activos deben coincidir con lo esperado")

	assert.Equal(t, len(expectedPeriods), len(dataGraph.ReactiveInductive))
	assert.Equal(t, len(expectedPeriods), len(dataGraph.ReactiveCapacitive))
	assert.Equal(t, len(expectedPeriods), len(dataGraph.Exported))
	for _, v := range dataGraph.ReactiveInductive {
		assert.Equal(t, 0.0, v)
	}
	for _, v := range dataGraph.ReactiveCapacitive {
		assert.Equal(t, 0.0, v)
	}
	for _, v := range dataGraph.Exported {
		assert.Equal(t, 0.0, v)
	}
}

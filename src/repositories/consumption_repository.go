package repositories

import (
	"time"

	"github.com/JuanConde27/energy-microservice/src/models"
	"gorm.io/gorm"
)

type ConsumptionRepository struct {
	DB *gorm.DB
}

func NewConsumptionRepository(db *gorm.DB) *ConsumptionRepository {
	return &ConsumptionRepository{DB: db}
}

func (r *ConsumptionRepository) GetMonthlyConsumption(meterIDs []int, startDate, endDate time.Time) ([]models.Consumption, error) {
	var consumptions []models.Consumption
	err := r.DB.Model(&models.Consumption{}).
		Select("meter_id, SUM(consumption) as consumption, DATE_TRUNC('month', timestamp) as period").
		Where("meter_id IN ? AND timestamp BETWEEN ? AND ?", meterIDs, startDate, endDate).
		Group("meter_id, period").
		Find(&consumptions).Error

	if err != nil {
		return nil, err
	}
	return consumptions, nil
}

func (r *ConsumptionRepository) GetWeeklyConsumption(meterIDs []int, startDate, endDate time.Time) ([]models.Consumption, error) {
	var consumptions []models.Consumption
	err := r.DB.Model(&models.Consumption{}).
		Select("meter_id, SUM(consumption) as consumption, DATE_TRUNC('week', timestamp) as period").
		Where("meter_id IN ? AND timestamp BETWEEN ? AND ?", meterIDs, startDate, endDate).
		Group("meter_id, period").
		Find(&consumptions).Error

	if err != nil {
		return nil, err
	}
	return consumptions, nil
}

func (r *ConsumptionRepository) GetDailyConsumption(meterIDs []int, startDate, endDate time.Time) ([]models.Consumption, error) {
	var consumptions []models.Consumption
	err := r.DB.Model(&models.Consumption{}).
		Select("meter_id, SUM(consumption) as consumption, DATE_TRUNC('day', timestamp) as period").
		Where("meter_id IN ? AND timestamp BETWEEN ? AND ?", meterIDs, startDate, endDate).
		Group("meter_id, period").
		Find(&consumptions).Error

	if err != nil {
		return nil, err
	}
	return consumptions, nil
}

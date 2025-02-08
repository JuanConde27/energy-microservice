package repositories

import (
	"github.com/JuanConde27/energy-microservice/src/models"
	"gorm.io/gorm"
	"time"
)

type ConsumptionRepository struct {
	DB *gorm.DB
}

func NewConsumptionRepository(db *gorm.DB) *ConsumptionRepository {
	return &ConsumptionRepository{DB: db}
}

func (r *ConsumptionRepository) GetMonthlyConsumption(meterIDs []int, startDate, endDate time.Time) ([]models.Consumption, error) {
	var consumptions []models.Consumption
	query := r.DB.Table("consumptions").
		Select("meter_id, SUM(consumption) as consumption, DATE_TRUNC('month', timestamp) as period").
		Where("meter_id IN (?) AND timestamp BETWEEN ? AND ?", meterIDs, startDate, endDate).
		Group("meter_id, period")

	if err := query.Find(&consumptions).Error; err != nil {
		return nil, err
	}
	return consumptions, nil
}

func (r *ConsumptionRepository) GetWeeklyConsumption(meterIDs []int, startDate, endDate time.Time) ([]models.Consumption, error) {
	var consumptions []models.Consumption
	query := r.DB.Table("consumptions").
		Select("meter_id, SUM(consumption) as consumption, DATE_TRUNC('week', timestamp) as period").
		Where("meter_id IN (?) AND timestamp BETWEEN ? AND ?", meterIDs, startDate, endDate).
		Group("meter_id, period")

	if err := query.Find(&consumptions).Error; err != nil {
		return nil, err
	}
	return consumptions, nil
}

func (r *ConsumptionRepository) GetDailyConsumption(meterIDs []int, startDate, endDate time.Time) ([]models.Consumption, error) {
	var consumptions []models.Consumption
	query := r.DB.Table("consumptions").
		Select("meter_id, SUM(consumption) as consumption, DATE_TRUNC('day', timestamp) as period").
		Where("meter_id IN (?) AND timestamp BETWEEN ? AND ?", meterIDs, startDate, endDate).
		Group("meter_id, period")

	if err := query.Find(&consumptions).Error; err != nil {
		return nil, err
	}
	return consumptions, nil
}

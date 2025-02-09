package repositories

import (
	"fmt"
	"github.com/lib/pq"
	"time"

	"github.com/JuanConde27/energy-microservice/src/models"
	"gorm.io/gorm"
)

type ConsumptionRepositoryInterface interface {
    GetConsumptionByPeriod(meterIDs []int, startDate, endDate time.Time, periodType string) ([]models.ConsumptionAggregate, error)
}

type ConsumptionRepository struct {
	DB *gorm.DB
}

func NewConsumptionRepository(db *gorm.DB) *ConsumptionRepository {
	return &ConsumptionRepository{DB: db}
}

func (r *ConsumptionRepository) GetConsumptionByPeriod(meterIDs []int, startDate, endDate time.Time, periodType string) ([]models.ConsumptionAggregate, error) {
	var consumptions []models.ConsumptionAggregate

	if periodType == "weekly" {
		query := `
            SELECT meter_id, 
                   COALESCE(SUM(consumption), 0) AS consumption,
                   ($1::timestamp + (((EXTRACT(EPOCH FROM timestamp - $1::timestamp))::int / (7 * 86400)) * interval '7 day')) AS period
            FROM consumptions
            WHERE meter_id = ANY($2)
              AND timestamp >= $1::timestamp AND timestamp < $3::timestamp
            GROUP BY meter_id, period
            ORDER BY period ASC;
        `
		err := r.DB.Raw(query, startDate, meterIDs, endDate).Scan(&consumptions).Error
		if err != nil {
			fmt.Println("❌ Error en consulta SQL:", err)
			return nil, err
		}

		return consumptions, nil
	}

	var dateTrunc string
	switch periodType {
	case "monthly":
		dateTrunc = "month"
	case "daily":
		dateTrunc = "day"
	default:
		dateTrunc = "day"
	}

	query := `
    SELECT meter_id, 
           COALESCE(SUM(consumption), 0) AS consumption, 
           DATE_TRUNC($1, timestamp) AS period
    FROM consumptions
    WHERE meter_id = ANY($2)
      AND timestamp >= $3 AND timestamp < ($4::timestamp + interval '1 day') 
    GROUP BY meter_id, period
    ORDER BY period ASC;
`
	err := r.DB.Raw(query, dateTrunc, pq.Array(meterIDs), startDate, endDate).Scan(&consumptions).Error

	if err != nil {
		fmt.Println("❌ Error en consulta SQL:", err)
		return nil, err
	}

	fmt.Println("✅ Datos obtenidos de la BD:")
	for _, c := range consumptions {
		fmt.Println("Meter ID:", c.MeterID, "Periodo:", c.Period, "Consumo:", c.Consumption)
	}

	return consumptions, nil
}

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

	fmt.Println("📌 Ejecutando GetConsumptionByPeriod")
	fmt.Println("🔹 Meter IDs:", meterIDs)
	fmt.Println("🔹 Start Date:", startDate)
	fmt.Println("🔹 End Date:", endDate)
	fmt.Println("🔹 Period Type:", periodType)

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
			fmt.Println("❌ Error en consulta SQL (weekly):", err)
			return nil, err
		}
		return consumptions, nil
	}

	if periodType == "monthly" {
		query := `
        SELECT 
            meters.meter_id,
            gs.month AS period,
            COALESCE(c.total_consumption, 0) AS consumption
        FROM (
            SELECT generate_series(date_trunc('month', $1::timestamp), date_trunc('month', $2::timestamp), '1 month') AS month
        ) gs
        CROSS JOIN (SELECT unnest($3::int[]) AS meter_id) meters
        LEFT JOIN (
            SELECT meter_id, date_trunc('month', timestamp) AS month, SUM(consumption) AS total_consumption
            FROM consumptions
            WHERE timestamp >= $1::timestamp AND timestamp < ($2::timestamp + interval '1 day')
            GROUP BY meter_id, date_trunc('month', timestamp)
        ) c ON c.meter_id = meters.meter_id AND c.month = gs.month
        ORDER BY gs.month ASC;
        `
		err := r.DB.Raw(query, startDate, endDate, meterIDs).Scan(&consumptions).Error
		if err != nil {
			fmt.Println("❌ Error en consulta SQL (monthly):", err)
			return nil, err
		}

		fmt.Println("✅ Datos obtenidos de la BD (Monthly):")
		for _, c := range consumptions {
			fmt.Println("Meter ID:", c.MeterID, "Periodo:", c.Period, "Consumo:", c.Consumption)
		}
		return consumptions, nil
	}

	var dateTrunc string
	switch periodType {
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
		fmt.Println("❌ Error en consulta SQL (daily):", err)
		return nil, err
	}
	fmt.Println("✅ Datos obtenidos de la BD (Daily):")
	for _, c := range consumptions {
		fmt.Println("Meter ID:", c.MeterID, "Periodo:", c.Period, "Consumo:", c.Consumption)
	}
	return consumptions, nil
}

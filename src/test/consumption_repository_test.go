package test

import (
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/stretchr/testify/assert"
	"github.com/JuanConde27/energy-microservice/src/repositories"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), 
	})
	return gormDB, mock
}

func TestGetConsumptionByPeriod(t *testing.T) {
	gormDB, mock := setupMockDB()
	repo := repositories.NewConsumptionRepository(gormDB)

	meterIDs := []int{1}
	startDate, _ := time.Parse("2006-01-02", "2023-06-01")
	endDate, _ := time.Parse("2006-01-02", "2023-06-09") 

	period1, _ := time.Parse("2006-01-02", "2023-06-01")
	period2, _ := time.Parse("2006-01-02", "2023-06-02")
	period3, _ := time.Parse("2006-01-02", "2023-06-03")

	mock.ExpectQuery(`SELECT meter_id, COALESCE\(SUM\(consumption\), 0\) AS consumption, DATE_TRUNC\(\$1, timestamp\) AS period FROM consumptions WHERE meter_id = ANY\(\$2\) AND timestamp >= \$3 AND timestamp < \(\$4::timestamp \+ interval '1 day'\) GROUP BY meter_id, period ORDER BY period ASC`).
		WithArgs("day", pq.Array(meterIDs), startDate, endDate). 
		WillReturnRows(sqlmock.NewRows([]string{"meter_id", "consumption", "period"}).
			AddRow(1, 139190.87165, period1).
			AddRow(1, 139653.67673, period2).
			AddRow(1, 140351.3282, period3),
		)

	consumptions, err := repo.GetConsumptionByPeriod(meterIDs, startDate, endDate, "daily")

	assert.NoError(t, err)
	assert.Len(t, consumptions, 3) 
	assert.Equal(t, float64(139190.87165), consumptions[0].Consumption)
}

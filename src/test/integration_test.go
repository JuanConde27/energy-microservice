package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/JuanConde27/energy-microservice/src/config"
	"github.com/JuanConde27/energy-microservice/src/server"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetConsumptionEndpoint(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err, "❌ Error al crear sqlmock")
    mock.MatchExpectationsInOrder(false) 
    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: db,
    }), &gorm.Config{})
    assert.NoError(t, err, "❌ Error al inicializar GORM")
    config.SetMockDB(gormDB)

    mock.ExpectQuery(`SELECT CURRENT_DATABASE\(\)`).
        WillReturnRows(sqlmock.NewRows([]string{"current_database"}).AddRow("energy_db"))
    mock.ExpectQuery(`SELECT count\(\*\) FROM information_schema\.tables`).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
    mock.ExpectQuery(`(?i)^SELECT .*FROM information_schema\.columns.*`).
        WithArgs("energy_db", "consumptions").
        WillReturnRows(sqlmock.NewRows([]string{"column_name"}).AddRow("id").AddRow("meter_id"))

    expectedRegex := `(?i)SELECT.*FROM consumptions.*WHERE.*`
    mock.ExpectQuery(expectedRegex).
        WithArgs(
            time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
            time.Date(2023, 6, 30, 23, 59, 59, 0, time.UTC),
            pq.Array([]int{1}),
        ).
        WillReturnRows(sqlmock.NewRows([]string{"meter_id", "consumption", "period"}).
            AddRow(1, 500.0, time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)).
            AddRow(1, 600.0, time.Date(2023, 6, 2, 0, 0, 0, 0, time.UTC)))

	os.Setenv("TEST_MODE", "true")
    router := server.SetupRouter()
    req, err := http.NewRequest("GET", "/consumption?meters_ids=1&start_date=2023-06-01&end_date=2023-06-30&kind_period=monthly", nil)
    assert.NoError(t, err)
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    fmt.Println("📌 Respuesta HTTP:", rr.Body.String())
    assert.Equal(t, http.StatusOK, rr.Code)

    var response map[string]interface{}
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Contains(t, response, "period")
    assert.Contains(t, response, "data_graph")

    dataGraph, ok := response["data_graph"].([]interface{})
    assert.True(t, ok)
    assert.Greater(t, len(dataGraph), 0)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "❌ Consultas mockeadas no cumplidas")
}
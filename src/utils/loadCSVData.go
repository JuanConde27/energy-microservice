package utils

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/JuanConde27/energy-microservice/src/config"
	"github.com/JuanConde27/energy-microservice/src/models"
)

func LoadCSVData(csvPath string) {
	db := config.GetConnection()
	defer config.CloseDb(db)

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Error abriendo el archivo CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error leyendo el archivo CSV: %v", err)
	}

	for _, row := range records {
		if len(row) < 4 {
			log.Println("Fila ignorada por datos incompletos:", row)
			continue
		}
	
		id := row[0]
		meterID, _ := strconv.Atoi(row[1])
		consumption, _ := strconv.ParseFloat(row[2], 64)
		timestamp, _ := time.Parse("2006-01-02 15:04:05-07", row[3])
	
		consumptionEntry := models.Consumption{
			ID:         id, 
			MeterID:    meterID,
			Consumption: consumption,
			Timestamp:  timestamp,
		}
	
		if err := db.Where("id = ?", id).FirstOrCreate(&consumptionEntry).Error; err != nil {
			log.Printf("Error insertando fila en la base de datos: %v", err)
		}
	}	

	log.Println("✅ Datos del CSV cargados exitosamente en la base de datos.")
}

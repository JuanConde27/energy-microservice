package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/JuanConde27/energy-microservice/src/models"
)

func GetConnection() *gorm.DB {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Println("⚠️  Advertencia: No se pudo cargar el archivo .env, usando variables de entorno del sistema.")
		}
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL no está definida en las variables de entorno.")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	log.Println("✅ Conexión exitosa a la base de datos.")
	return database
}

func CloseDb(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
}

func Migrate() {
	db := GetConnection()
	defer CloseDb(db)

	db.AutoMigrate(&models.Consumption{})
	log.Println("✅ Migración completada.")
}

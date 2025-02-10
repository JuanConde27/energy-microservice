package config

import (
    "log"
    "os"
    "sync"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/JuanConde27/energy-microservice/src/models"
)

var (
    mockDB   *gorm.DB
    isMocked bool
    mu       sync.Mutex
)

func SetMockDB(mock *gorm.DB) {
    mu.Lock()
    defer mu.Unlock()
    mockDB = mock
    isMocked = true
    log.Println("✅ Base de datos mockeada asignada correctamente.")
}

func GetMockDB() (*gorm.DB, bool) {
    mu.Lock()
    defer mu.Unlock()
    return mockDB, isMocked
}

func GetConnection() *gorm.DB {
    if db, mocked := GetMockDB(); mocked {
        log.Println("✅ Usando base de datos mockeada.")
        return db
    }

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

    log.Println("✅ Conexión exitosa a la base de datos real.")
    return database
}

func CloseDb(db *gorm.DB) {
    if _, mocked := GetMockDB(); mocked {
        log.Println("✅ Cierre de base de datos omitido porque se está usando una base de datos mockeada.")
        return
    }

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
	log.Println("✅ Migraciones de la base de datos completadas.")
}
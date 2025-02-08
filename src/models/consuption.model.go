package models

import (
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Consumption struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	MeterID    int       `gorm:"not null" json:"meter_id"`
	Consumption float64  `gorm:"not null" json:"consumption"`
	Timestamp  time.Time `gorm:"not null" json:"timestamp"`
}

func (c *Consumption) TableName() string {
	return "consumptions"
}

func (c *Consumption) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}


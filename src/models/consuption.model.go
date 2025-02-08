package models

import (
	"time"
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


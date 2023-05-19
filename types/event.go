package types

import (
	"time"

	"gorm.io/gorm"
)

type EventType int

const (
	PartyEventType = iota
	FestivalEventType
	ConcertEventType
)

type Event struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	OrganizationID uint
	Name           string    `json:"name"`
	Type           EventType `json:"type"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	CountryCode    string    `json:"countryCode"`
	City           string    `json:"city"`
}

package types

import (
	"fmt"
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

type CreateEventParams struct {
	Name        string    `json:"name"`
	Type        EventType `json:"type"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	CountryCode string    `json:"countryCode"`
	City        string    `json:"city"`
}

type UpdateEventParams struct {
	Name        string    `json:"name"`
	Type        EventType `json:"type"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	CountryCode string    `json:"countryCode"`
	City        string    `json:"city"`
}

func (p CreateEventParams) Validate() map[string]string {
	errors := map[string]string{}

	if len(p.Name) < minNameLen {
		errors["name"] = fmt.Sprintf("name length should be at least %d", minNameLen)
	}

	return errors
}

func (p UpdateEventParams) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	if len(p.Name) > 0 {
		m["name"] = p.Name
	}
	return m
}

func NewEventFromParams(params CreateEventParams) (*Event, error) {
	return &Event{
		Name:        params.Name,
		Type:        params.Type,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
		CountryCode: params.CountryCode,
		City:        params.City,
	}, nil
}

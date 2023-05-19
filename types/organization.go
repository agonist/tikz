package types

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	minNameLen = 3
)

type Organization struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name   string  `json:"name"`
	Events []Event `json:"events"`
}

type CreateOrgParams struct {
	Name string `json:"name"`
}

func (p CreateOrgParams) Validate() map[string]string {
	errors := map[string]string{}

	if len(p.Name) < minNameLen {
		errors["name"] = fmt.Sprintf("name length should be at least %d", minNameLen)
	}

	return errors
}

type UpdateOrgParams struct {
	Name string `json:"name"`
}

func (p UpdateOrgParams) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	if len(p.Name) > 0 {
		m["name"] = p.Name
	}
	return m
}

func NewOrgFromParams(params CreateOrgParams) (*Organization, error) {
	return &Organization{
		Name: params.Name,
	}, nil
}

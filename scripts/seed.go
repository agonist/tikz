package main

import (
	"fmt"
	"time"

	"github.com/agonist/hotel-reservation/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dburi = "host=localhost user=admin password=supersecret dbname=ticketing port=5432 sslmode=disable"

func main() {
	fmt.Println("seeding the database...")

	db, err := gorm.Open(postgres.Open(dburi), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the db")
	}
	err = db.AutoMigrate(&types.User{})
	err = db.AutoMigrate(&types.Organization{})
	err = db.AutoMigrate(&types.Event{})

	if err != nil {
		panic("failed ti run migrations")
	}

	organization1 := types.Organization{
		Name: "Modem",
		Events: []types.Event{
			{
				Name:        "Modem Festival London Teaser",
				Type:        types.PartyEventType,
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(time.Hour * 24),
				CountryCode: "UK",
				City:        "London",
			},
			{
				Name:        "Modem Festival",
				Type:        types.FestivalEventType,
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(3 * time.Hour * 24),
				CountryCode: "CR",
				City:        "Slunj",
			},
		},
	}
	organization2 := types.Organization{
		Name: "Hellfest",
		Events: []types.Event{
			{
				Name:        "Hellfest Festival",
				Type:        types.FestivalEventType,
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(time.Hour * 24),
				CountryCode: "FR",
				City:        "Clisson",
			},
		},
	}

	organization3 := types.Organization{
		Name: "MoP",
		Events: []types.Event{
			{
				Name:        "Master Of Puppets Festival Teaser",
				Type:        types.PartyEventType,
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(time.Hour * 24),
				CountryCode: "UK",
				City:        "London",
			},
			{
				Name:        "Master Of Puppets Festival",
				Type:        types.FestivalEventType,
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(7 * time.Hour * 24),
				CountryCode: "CZ",
				City:        "SomeTownName",
			},
		},
	}

	organizations := []*types.Organization{
		&organization1,
		&organization2,
		&organization3,
	}
	res := db.Create(organizations)
	if res.Error != nil {
		fmt.Errorf(res.Error.Error())
	}
}

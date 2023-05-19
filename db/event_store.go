package db

import (
	"github.com/agonist/hotel-reservation/types"
	"gorm.io/gorm"
)

// just using org instead of organization cause its fucking long

type EventStore interface {
	Dropper

	GetByID(int) (*types.Event, error)
	GetAll() (*[]types.Event, error)
	Insert(*types.Event) (*types.Event, error)
	Delete(int) error
	Update(eventID int, updated types.UpdateEventParams) error
}

type PgEventStore struct {
	client *gorm.DB
}

func NewPgEventStore(c *gorm.DB) *PgEventStore {
	return &PgEventStore{
		client: c,
	}
}

func (s *PgEventStore) Drop() error {
	s.client.Migrator().DropTable(&types.Event{})
	return nil
}

func (s *PgEventStore) Insert(event *types.Event) (*types.Event, error) {
	res := s.client.Create(&event)
	if res.Error != nil {
		return nil, res.Error
	}

	return event, nil
}

func (s *PgEventStore) GetByID(id int) (*types.Event, error) {
	var event types.Event

	res := s.client.First(&event, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &event, nil
}

func (s *PgEventStore) GetAll() (*[]types.Event, error) {
	var events []types.Event

	res := s.client.Find(&events)
	if res.Error != nil {
		return nil, res.Error
	}

	return &events, nil
}

func (s *PgEventStore) Update(id int, update types.UpdateEventParams) error {
	res := s.client.Model(&types.Event{ID: uint(id)}).Updates(update.ToMap())

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *PgEventStore) Delete(id int) error {
	res := s.client.Delete(&types.Event{}, id)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

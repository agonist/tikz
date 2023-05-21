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
	Insert(event *types.Event) (*types.Event, error)
	Delete(int) error
	Update(eventID int, updated types.UpdateEventParams) error
}

type PgEventStore struct {
	client   *gorm.DB
	orgStore OrgStore
}

func NewPgEventStore(c *gorm.DB, orgStore OrgStore) *PgEventStore {
	return &PgEventStore{
		client:   c,
		orgStore: orgStore,
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
	org, err := s.orgStore.GetByID(int(event.OrganizationID))
	if err != nil {
		return nil, err
	}
	org.Events = append(org.Events, *event)
	s.client.Save(org)
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

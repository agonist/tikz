package db

import (
	"github.com/agonist/hotel-reservation/types"
	"gorm.io/gorm"
)

// just using org instead of organization cause its fucking long

type OrgStore interface {
	Dropper

	GetByID(int) (*types.Organization, error)
	GetAll() (*[]types.Organization, error)
	Insert(*types.Organization) (*types.Organization, error)
	Delete(int) error
	Update(orgId int, updated types.UpdateOrgParams) error
}

type PgOrgStore struct {
	client *gorm.DB
}

func NewPgOrgStore(c *gorm.DB) *PgOrgStore {
	return &PgOrgStore{
		client: c,
	}
}

func (s *PgOrgStore) Drop() error {
	s.client.Migrator().DropTable(&types.Organization{})
	return nil
}

func (s *PgOrgStore) Insert(org *types.Organization) (*types.Organization, error) {
	res := s.client.Create(&org)
	if res.Error != nil {
		return nil, res.Error
	}

	return org, nil
}

func (s *PgOrgStore) GetByID(id int) (*types.Organization, error) {
	var org types.Organization

	res := s.client.First(&org, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &org, nil
}

func (s *PgOrgStore) GetAll() (*[]types.Organization, error) {
	var orgs []types.Organization

	res := s.client.Find(&orgs)
	if res.Error != nil {
		return nil, res.Error
	}

	return &orgs, nil
}

func (s *PgOrgStore) Update(id int, update types.UpdateOrgParams) error {
	res := s.client.Model(&types.Organization{ID: uint(id)}).Updates(update.ToMap())

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *PgOrgStore) Delete(id int) error {
	res := s.client.Delete(&types.Organization{}, id)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

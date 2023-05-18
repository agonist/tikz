package db

import (
	"github.com/agonist/hotel-reservation/types"
	"gorm.io/gorm"
)

const userColl = "users"

type Dropper interface {
	Drop() error
}

type UserStore interface {
	Dropper

	GetUserByID(int) (*types.User, error)
	GetUsers() (*[]types.User, error)
	InsertUser(*types.User) (*types.User, error)
	DeleteUser(int) error
	UpdateUser(userId int, update types.UpdateUserParams) error
}

type PgUserStore struct {
	client *gorm.DB
}

func NewPgUserStore(c *gorm.DB) *PgUserStore {
	return &PgUserStore{
		client: c,
	}
}

func (s *PgUserStore) Drop() error {
	s.client.Migrator().DropTable(&types.User{})
	return nil
}

func (s *PgUserStore) InsertUser(user *types.User) (*types.User, error) {
	res := s.client.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (s *PgUserStore) GetUserByID(id int) (*types.User, error) {
	var user types.User

	res := s.client.First(&user, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (s *PgUserStore) GetUsers() (*[]types.User, error) {
	var users []types.User

	res := s.client.Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return &users, nil
}

func (s *PgUserStore) UpdateUser(userId int, params types.UpdateUserParams) error {
	res := s.client.Model(&types.User{ID: uint(userId)}).Updates(params.ToMap())

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *PgUserStore) DeleteUser(id int) error {
	res := s.client.Delete(&types.User{}, id)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

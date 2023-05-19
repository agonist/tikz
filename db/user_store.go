package db

import (
	"errors"

	"github.com/agonist/hotel-reservation/types"
	"gorm.io/gorm"
)

const userColl = "users"

type UserStore interface {
	Dropper

	GetByID(int) (*types.User, error)
	GetAll() (*[]types.User, error)
	Insert(*types.User) (*types.User, error)
	Delete(int) error
	Update(userId int, update types.UpdateUserParams) error
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

func (s *PgUserStore) Insert(user *types.User) (*types.User, error) {
	res := s.client.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (s *PgUserStore) GetByID(id int) (*types.User, error) {
	var user types.User

	if err := s.client.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return &user, nil
}

func (s *PgUserStore) GetAll() (*[]types.User, error) {
	var users []types.User

	res := s.client.Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return &users, nil
}

func (s *PgUserStore) Update(userId int, params types.UpdateUserParams) error {
	res := s.client.Model(&types.User{ID: uint(userId)}).Updates(params.ToMap())

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *PgUserStore) Delete(id int) error {
	res := s.client.Delete(&types.User{}, id)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

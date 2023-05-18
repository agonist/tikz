package types

import (
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"-"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p UpdateUserParams) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	if len(p.FirstName) > 0 {
		m["first_name"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["last_name"] = p.LastName
	}
	return m
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}

	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstname length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastname length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email is invalid")
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}

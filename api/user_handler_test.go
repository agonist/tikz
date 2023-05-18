package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/agonist/hotel-reservation/db"
	"github.com/agonist/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dburi = "host=localhost user=admin password=supersecret dbname=ticketing port=5432 sslmode=disable"

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	tdb, err := gorm.Open(postgres.Open(dburi), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the db")
	}
	err = tdb.AutoMigrate(&types.User{})
	if err != nil {
		panic("failed ti run migrations")
	}

	return &testdb{
		UserStore: db.NewPgUserStore(tdb),
	}
}

var (
	user1 = types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "John",
		LastName:  "Bob",
		Password:  "123456789",
	}
	user2 = types.CreateUserParams{
		Email:     "foo@bar.com",
		FirstName: "Alice",
		LastName:  "Bread",
		Password:  "123456789",
	}
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	b, _ := json.Marshal(user1)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	assert.NotZero(t, user.ID, "ID should not be 0")
	assert.Equal(t, user.Email, user1.Email, "email should be equal")
	assert.Equal(t, user.FirstName, user1.FirstName, "firstName should be equal")
	assert.Equal(t, user.LastName, user1.LastName, "lastName should be equal")
}

func TestListUSer(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)
	app.Get("/", userHandler.HandleListUsers)

	b, _ := json.Marshal(user1)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	_, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	b, _ = json.Marshal(user2)
	req = httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	_, err = app.Test(req)
	if err != nil {
		t.Error(err)
	}

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Error(err)
	}

	var user []types.User
	json.NewDecoder(resp.Body).Decode(&user)
	assert.Len(t, user, 2)
}

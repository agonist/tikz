package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/agonist/hotel-reservation/db"
	"github.com/agonist/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel_reservation_test"

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
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

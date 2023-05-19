package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agonist/hotel-reservation/types"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/stretchr/testify/assert"
)

var (
	userInvalid = types.CreateUserParams{
		Email:     "somefoo.com",
		FirstName: "J",
		LastName:  "B",
		Password:  "12",
	}
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

func postUser(user types.CreateUserParams) *http.Request {
	b, _ := sonic.Marshal(user)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	return req
}

func listAllUser() *http.Request {
	req := httptest.NewRequest("GET", "/api/v1/user", nil)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func getUserByID(id int) *http.Request {
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/user/%d", id), nil)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func deleteUser(id int) *http.Request {
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/user/%d", id), nil)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func TestPostUser(t *testing.T) {
	tc := setup(t)
	defer tc.teardown(t)

	resp, err := tc.app.Test(postUser(userInvalid))
	if err != nil {
		t.Error(err)
	}

	resp, err = tc.app.Test(postUser(user1))
	if err != nil {
		t.Error(err)
	}

	var user types.User
	decoder.NewStreamDecoder(resp.Body).Decode(&user)
	assert.NotZero(t, user.ID, "ID should not be 0")
	assert.Equal(t, user.Email, user1.Email, "email should be equal")
	assert.Equal(t, user.FirstName, user1.FirstName, "firstName should be equal")
	assert.Equal(t, user.LastName, user1.LastName, "lastName should be equal")
}

func TestListUser(t *testing.T) {
	tc := setup(t)
	defer tc.teardown(t)

	resp, err := tc.app.Test(postUser(user1))
	if err != nil {
		t.Error(err)
	}

	resp, err = tc.app.Test(postUser(user2))
	if err != nil {
		t.Error(err)
	}

	resp, err = tc.app.Test(listAllUser())

	if err != nil {
		t.Error(err)
	}

	var users []types.User
	decoder.NewStreamDecoder(resp.Body).Decode(&users)
	assert.Len(t, users, 2)
}

func TestGetUserById(t *testing.T) {
	tc := setup(t)
	defer tc.teardown(t)

	_, err := tc.app.Test(postUser(user1))
	if err != nil {
		t.Error(err)
	}

	_, err = tc.app.Test(postUser(user2))
	if err != nil {
		t.Error(err)
	}

	resp, err := tc.app.Test(getUserByID(1))
	if err != nil {
		t.Error(err)
	}

	var user types.User
	decoder.NewStreamDecoder(resp.Body).Decode(&user)
	assert.NotZero(t, user.ID, "ID should not be 0")
	assert.Equal(t, user.Email, user1.Email, "email should be equal")
	assert.Equal(t, user.FirstName, user1.FirstName, "firstName should be equal")
	assert.Equal(t, user.LastName, user1.LastName, "lastName should be equal")
}

func TestDeleteUser(t *testing.T) {
	tc := setup(t)
	defer tc.teardown(t)

	_, err := tc.app.Test(postUser(user1))
	if err != nil {
		t.Error(err)
	}

	_, err = tc.app.Test(deleteUser(1))
	if err != nil {
		t.Error(err)
	}

	resp, err := tc.app.Test(getUserByID(1))
	if err != nil {
		t.Error(err)
	}
	var e terr
	decoder.NewStreamDecoder(resp.Body).Decode(&e)
	assert.Equal(t, e.Msg, "User not found")
}

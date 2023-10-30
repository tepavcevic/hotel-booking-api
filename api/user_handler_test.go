package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/types"
)

func TestCreateUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandleCreateUser)

	params := types.CreateUserParams{
		FirstName: "some",
		LastName:  "user",
		Email:     "some@user.com",
		Password:  "somepassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Error("expected user ID to be set")
	}
	if len(user.PasswordHash) > 0 {
		t.Error("expected user password hash to be ommited")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected %s, got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected %s, got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected %s, got %s", params.Email, user.Email)
	}
}

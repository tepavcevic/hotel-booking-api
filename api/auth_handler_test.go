package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/db/fixtures"
)

func TestAuthenticate(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := fixtures.AddUser(tdb.Store, "Mario", "Dreznjak", true)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "mario@dreznjak.com",
		Password: "greatestpasswordever",
	}
	b, _ := json.Marshal(authParams)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 but got %v", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if len(authResp.Token) == 0 {
		t.Fatalf("expected JWT token to be present in the response")
	}
	// set password hash to "" because we don't return it in any json response
	insertedUser.PasswordHash = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatal("expected the user to be inserted user")
	}
}

func TestAuthenticatePasswordFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	fixtures.AddUser(tdb.Store, "Mario", "Dreznjak", true)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "mario@dreznjak.com",
		Password: "wrongpsswrd",
	}
	b, _ := json.Marshal(authParams)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 but got %v", resp.StatusCode)
	}
	var genResp genericResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatal("expected genResp Type to be error, got:", genResp.Type)
	}
	if genResp.Message != "invalid credentials" {
		t.Fatal("expected messagee to be invalid credentials, got:", genResp.Message)
	}
}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/api/middleware"
	"github.com/tepavcevic/hotel-reservation/db/fixtures"
	"github.com/tepavcevic/hotel-reservation/types"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user      = fixtures.AddUser(tdb.Store, "James", "Rodriguez", false)
		hotel     = fixtures.AddHotel(tdb.Store, "Hilton", "USA", 9, nil)
		room      = fixtures.AddRoom(tdb.Store, types.KingSizeRoom, true, 1799.99, hotel.ID)
		booking   = fixtures.AddBooking(
			tdb.Store,
			room.ID,
			user.ID,
			time.Now(),
			time.Now().AddDate(0, 0, 12),
			2,
		)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(tdb.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 response, got %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s, got %s", booking.ID, have.ID)
	}
	fmt.Println(bookings)

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized response, got %d", resp.StatusCode)
	}
}

func TestUserGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user        = fixtures.AddUser(tdb.Store, "James", "Rodriguez", false)
		nonAuthUser = fixtures.AddUser(tdb.Store, "Michael", "Ballack", false)
		hotel       = fixtures.AddHotel(tdb.Store, "Hilton", "USA", 9, nil)
		room        = fixtures.AddRoom(tdb.Store, types.KingSizeRoom, true, 1799.99, hotel.ID)
		booking     = fixtures.AddBooking(
			tdb.Store,
			room.ID,
			user.ID,
			time.Now(),
			time.Now().AddDate(0, 0, 12),
			2,
		)
		app            = fiber.New()
		route          = app.Group("/:id", middleware.JWTAuthentication(tdb.User))
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	_ = booking
	route.Get("/", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("Authorization", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 response, got %d", resp.StatusCode)
	}
	var bookingResponse *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResponse); err != nil {
		t.Fatal(err)
	}
	if bookingResponse.ID != booking.ID {
		t.Fatalf("expected %s, got %s", booking.ID, bookingResponse.ID)
	}
	if bookingResponse.UserID != booking.UserID {
		t.Fatalf("expected %s, got %s", booking.UserID, bookingResponse.UserID)
	}
	fmt.Println(bookingResponse)

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("Authorization", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("want non 200 response, got %d", resp.StatusCode)
	}
}

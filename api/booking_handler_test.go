package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/tepavcevic/hotel-reservation/db/fixtures"
	"github.com/tepavcevic/hotel-reservation/types"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.Store, "James", "Rodriguez", false)
	hotel := fixtures.AddHotel(tdb.Store, "Hilton", "USA", 9, nil)
	room := fixtures.AddRoom(tdb.Store, types.KingSizeRoom, true, 1799.99, hotel.ID)
	booking := fixtures.AddBooking(
		tdb.Store,
		room.ID,
		user.ID,
		time.Now(),
		time.Now().AddDate(0, 0, 12),
		2,
	)
	fmt.Println("booking", booking)
}

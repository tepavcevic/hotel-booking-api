package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tepavcevic/hotel-reservation/api"
	"github.com/tepavcevic/hotel-reservation/db"
	"github.com/tepavcevic/hotel-reservation/db/fixtures"
	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	hotelCollName   = "hotels"
	roomCollName    = "rooms"
	userCollName    = "users"
	bookingCollName = "bookings"
)

var (
	hotelStore db.HotelStore
)

func main() {
	ctx := context.Background()
	var err error
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client, db.DBNAME, hotelCollName)
	store := db.Store{
		User:    db.NewMongoUserStore(client, db.DBNAME, userCollName),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, db.DBNAME, roomCollName, hotelStore),
		Booking: db.NewMongoBookingStore(client, db.DBNAME, bookingCollName),
	}
	user := fixtures.AddUser(&store, "Jaes", "Nunes", false)
	fmt.Println("Jaes ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(&store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(&store, "Regent", "Austria", 8, nil)
	room := fixtures.AddRoom(&store, types.NormalRoom, false, 123.33, hotel.ID)
	booking := fixtures.AddBooking(&store, room.ID, user.ID, time.Now(), time.Now().AddDate(0, 0, 5), 3)
	fmt.Println(booking)
}

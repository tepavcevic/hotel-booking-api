package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		dbURI  = os.Getenv("DB_URI")
		dbName = os.Getenv("DB_NAME")
		ctx    = context.Background()
		err    error
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(dbName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client, dbName, hotelCollName)
	store := db.Store{
		User:    db.NewMongoUserStore(client, dbName, userCollName),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, dbName, roomCollName, hotelStore),
		Booking: db.NewMongoBookingStore(client, dbName, bookingCollName),
	}
	user := fixtures.AddUser(&store, "Jaes", "Nunes", false)
	fmt.Println("Jaes ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(&store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(&store, "Regent", "Austria", 8, nil)
	room := fixtures.AddRoom(&store, types.NormalRoom, false, 123.33, hotel.ID)
	booking := fixtures.AddBooking(&store, room.ID, user.ID, time.Now(), time.Now().AddDate(0, 0, 5), 3)
	fmt.Println(booking)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random fake hotel %d", i)
		location := fmt.Sprintf("random location %d", i)
		fixtures.AddHotel(&store, name, location, i%10+1, nil)
	}
}

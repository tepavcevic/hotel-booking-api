package main

import (
	"context"
	"log"
	"time"

	"github.com/tepavcevic/hotel-reservation/db"
	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, email, firstName, lastName, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	})
	user.IsAdmin = isAdmin
	if err != nil {
		log.Fatal(err)
	}
	dbUser, err := userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	return dbUser
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	dbHotel, err := hotelStore.Create(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return dbHotel
}

func seedRoom(size types.Size, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		HotelID: hotelID,
		Size:    size,
		Seaside: seaside,
		Price:   price,
	}
	dbRoom, err := roomStore.CreateRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}
	return dbRoom
}

func seedBooking(roomID, userID primitive.ObjectID, fromDate, tillDate time.Time, numPersons int, cancelled bool) {
	booking := types.Booking{
		RoomID:     roomID,
		UserID:     userID,
		FromDate:   fromDate,
		TillDate:   tillDate,
		NumPersons: numPersons,
		Cancelled:  cancelled,
	}
	_, err := bookingStore.CreateBooking(ctx, &booking)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	seedHotel("Transilvania", "Romania", 8)
	seedHotel("Moskva", "Serbia", 7)
	swissotel := seedHotel("Swissotel", "Bosnia", 8)
	json := seedUser(false, "son@momoaa.com", "Momo", "Jsoon", "greatestpasswordever")
	seedUser(true, "mario@dreznjak.com", "Mario", "Dreznjak", "admin123")
	room := seedRoom(types.NormalRoom, false, 87.99, swissotel.ID)
	seedBooking(room.ID, json.ID, time.Now(), time.Now().AddDate(0, 0, 4), 3, false)
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client, db.DBNAME, hotelCollName)
	roomStore = db.NewMongoRoomStore(client, db.DBNAME, roomCollName, hotelStore)
	userStore = db.NewMongoUserStore(client, db.DBNAME, userCollName)
	bookingStore = db.NewMongoBookingStore(client, db.DBNAME, bookingCollName)
	if err := client.Database(db.DBNAME).Collection(hotelCollName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Collection(roomCollName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Collection(userCollName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Collection(bookingCollName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
}

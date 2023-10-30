package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tepavcevic/hotel-reservation/db"
	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	user.IsAdmin = admin
	if err != nil {
		log.Fatal(err)
	}
	dbUser, err := store.User.CreateUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return dbUser
}

func AddHotel(store *db.Store, name, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}
	dbHotel, err := store.Hotel.Create(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return dbHotel
}

func AddRoom(store *db.Store, size types.Size, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		HotelID: hotelID,
		Size:    size,
		Seaside: ss,
		Price:   price,
	}
	dbRoom, err := store.Room.CreateRoom(context.Background(), &room)
	if err != nil {
		log.Fatal(err)
	}
	return dbRoom
}

func AddBooking(store *db.Store, roomID, userID primitive.ObjectID, from, till time.Time, numP int) *types.Booking {
	booking := types.Booking{
		RoomID:     roomID,
		UserID:     userID,
		FromDate:   from,
		TillDate:   till,
		NumPersons: numP,
		Cancelled:  false,
	}
	dbBooking, err := store.Booking.CreateBooking(context.Background(), &booking)
	if err != nil {
		log.Fatal(err)
	}
	return dbBooking
}

package main

import (
	"context"
	"log"

	"github.com/tepavcevic/hotel-reservation/db"
	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	hotelCollName = "hotels"
	roomCollName  = "rooms"
	userCollName  = "users"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(isAdmin bool, email, firstName, lastName, password string) {
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
	_, err = userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Size:    types.SmallRoom,
			Seaside: false,
			Price:   65,
		},
		{
			Size:    types.NormalRoom,
			Seaside: true,
			Price:   99,
		},
		{
			Size:    types.KingSizeRoom,
			Seaside: true,
			Price:   189,
		},
	}
	dbHotel, err := hotelStore.Create(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = dbHotel.ID
		dbRoom, err := roomStore.CreateRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		dbHotel.Rooms = append(dbHotel.Rooms, dbRoom.ID)
	}
}

func main() {
	seedHotel("Transilvania", "Romania", 8)
	seedHotel("Moskva", "Serbia", 7)
	seedHotel("Swissotel", "Bosnia", 8)
	seedUser(false, "son@momoaa.com", "Momo", "Jsoon", "greatestpasswordever")
	seedUser(true, "mario@dreznjak.com", "Mario", "Dreznjak", "admin123")
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
	if err := client.Database(db.DBNAME).Collection(hotelCollName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Collection(roomCollName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Collection(userCollName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
}

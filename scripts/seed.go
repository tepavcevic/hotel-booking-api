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
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Type:      types.DoubleRoomType,
			BasePrice: 65,
		},
		{
			Type:      types.SingleRoomType,
			BasePrice: 55,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 165,
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
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client, db.DBNAME, hotelCollName)
	roomStore = db.NewMongoRoomStore(client, db.DBNAME, roomCollName, hotelStore)
	if err := client.Database(db.DBNAME).Collection(hotelCollName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Collection(roomCollName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
}

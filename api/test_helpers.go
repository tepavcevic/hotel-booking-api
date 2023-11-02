package api

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/tepavcevic/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDBName          = "test-hotel-reservation"
	testUserCollName    = "test-users"
	testHotelCollName   = "test-hotels"
	testRoomCollName    = "test-rooms"
	testBookingCollName = "test-bookings"
)

type testdb struct {
	*db.Store
	client *mongo.Client
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal()
	}
	dbURI := os.Getenv("DB_URI_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, testDBName, testHotelCollName)
	return &testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client, testDBName, testUserCollName),
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, testDBName, testRoomCollName, hotelStore),
			Booking: db.NewMongoBookingStore(client, testDBName, testBookingCollName),
		},
	}
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(testDBName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

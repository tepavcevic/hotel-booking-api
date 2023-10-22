package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/api"
	"github.com/tepavcevic/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	usersColl  = "users"
	hotelsColl = "hotels"
	roomsColl  = "rooms"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":8080", "Port for our server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	// store, handler and fiber initializations
	var (
		userStore    = db.NewMongoUserStore(client, db.DBNAME, usersColl)
		hotelStore   = db.NewMongoHotelStore(client, db.DBNAME, hotelsColl)
		roomsStore   = db.NewMongoRoomStore(client, db.DBNAME, roomsColl, hotelStore)
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(hotelStore, roomsStore)
		app          = fiber.New(config)
		apiv1        = app.Group("/api/v1")
	)

	apiv1.Post("users", userHandler.HandleCreateUser)
	apiv1.Get("/users/:id", userHandler.HandleGetUserById)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Put("/users/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)

	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	app.Listen(*listenAddr)
}

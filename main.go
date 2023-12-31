package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tepavcevic/hotel-reservation/api"
	"github.com/tepavcevic/hotel-reservation/api/middleware"
	"github.com/tepavcevic/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	usersColl    = "users"
	hotelsColl   = "hotels"
	roomsColl    = "rooms"
	bookingsColl = "bookings"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	var (
		dbURI  = os.Getenv("DB_URI")
		dbName = os.Getenv("DB_NAME")
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}
	// store, handler and fiber initializations
	var (
		userStore    = db.NewMongoUserStore(client, dbName, usersColl)
		hotelStore   = db.NewMongoHotelStore(client, dbName, hotelsColl)
		roomsStore   = db.NewMongoRoomStore(client, dbName, roomsColl, hotelStore)
		bookingStore = db.NewMongoBookingStore(client, dbName, bookingsColl)
		store        = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomsStore,
			Booking: bookingStore,
		}
		authHandler    = api.NewAuthHandler(userStore)
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", middleware.AdminAuth)
	)

	auth.Post("/auth", authHandler.HandleAuthenticate)

	apiv1.Post("users", userHandler.HandleCreateUser)
	apiv1.Get("/users/:id", userHandler.HandleGetUserById)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Put("/users/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)

	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id", hotelHandler.HandleGetHotelByID)
	apiv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)

	apiv1.Post("/rooms/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/rooms", roomHandler.HandleGetRooms)

	apiv1.Get("/bookings/:id", bookingHandler.HandleGetBooking)
	apiv1.Put("/bookings/:id", bookingHandler.HandleCancelBooking)

	admin.Get("/bookings", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

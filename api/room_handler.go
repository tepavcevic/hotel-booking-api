package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/db"
	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("invalid booking date")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (rh *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := rh.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(http.StatusNotFound).JSON(genericResponse{
				Type:    "error",
				Message: "no rooms found",
			})
		}
		return err
	}
	return c.JSON(rooms)
}

func (rh *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return c.Status(http.StatusBadRequest).JSON(genericResponse{
			Type:    "error",
			Message: "bad request",
		})
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomOID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type:    "error",
			Message: "internal server error",
		})
	}
	ok, err = rh.isRoomAvailableForBooking(c.Context(), roomOID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResponse{
			Type:    "error",
			Message: fmt.Sprintf("room %s is already booked for this period", c.Params("id")),
		})
	}
	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomOID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}
	dbbooking, err := rh.store.Booking.CreateBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", booking)
	return c.JSON(dbbooking)
}

func (rh *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomOID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomOID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := rh.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	if len(bookings) > 0 {
		return false, nil
	}
	return true, nil
}

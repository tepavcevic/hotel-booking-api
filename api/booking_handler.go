package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (bh *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := bh.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (bh *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := bh.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if user.ID != booking.UserID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type:    "error",
			Message: "unauthorized",
		})
	}
	return c.JSON(booking)
}

func (bh *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	bookingID := c.Params("id")
	booking, err := bh.store.Booking.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type:    "error",
			Message: "unauthorized",
		})
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type:    "error",
			Message: "unauthorized",
		})
	}
	if err := bh.store.Booking.UpdateBooking(c.Context(), bookingID, bson.M{"cancelled": true}); err != nil {
		return err
	}
	return c.JSON(genericResponse{
		Type:    "message",
		Message: "booking cancelled",
	})
}

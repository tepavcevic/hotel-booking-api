package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/db"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := db.Map{"hotelID": id}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	filter := db.Map{}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

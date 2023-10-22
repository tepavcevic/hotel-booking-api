package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	filter := bson.M{}
	hotels, err := h.hotelStore.GetHotels(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

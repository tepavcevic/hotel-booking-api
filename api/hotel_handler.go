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

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

type HotelQueryParams struct {
	db.PaginationParams
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return fiber.ErrBadRequest
	}
	filter := db.Map{
		"rating": qparams.Rating,
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &qparams.PaginationParams)
	if err != nil {
		return err
	}
	resp := ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(qparams.Page),
	}
	return c.JSON(resp)
}

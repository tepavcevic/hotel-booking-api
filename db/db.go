package db

const (
	DBNAME = "hotel-reservation"
	DBURI  = "mongodb://localhost:27017"
)

type PaginationParams struct {
	Page  int64
	Limit int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int                  `bson:"rating" json:"rating"`
}

type Size string

const (
	SmallRoom    = "small"
	NormalRoom   = "normal"
	KingSizeRoom = "kingsize"
)

type Room struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// SmallRoom, NormalRoom, KingSizeRoom
	Size    Size               `bson:"size" json:"size"`
	Seaside bool               `bson:"seaside" json:"seaside"`
	Price   float64            `bson:"price" json:"price"`
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}

package db

import (
	"context"

	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	CreateHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, bson.M, bson.M) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(c *mongo.Client, dbname, collName string) *MongoHotelStore {
	return &MongoHotelStore{
		client: c,
		coll:   c.Database(dbname).Collection(collName),
	}
}

func (store *MongoHotelStore) CreateHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := store.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (store *MongoHotelStore) Update(ctx context.Context, filter, update bson.M) error {
	_, err := store.coll.UpdateOne(ctx, filter, update)
	return err
}

package db

import (
	"context"

	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	GetHotels(context.Context, bson.M) ([]*types.Hotel, error)
	Create(context.Context, *types.Hotel) (*types.Hotel, error)
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

func (store *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	var hotels []*types.Hotel
	cur, err := store.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (store *MongoHotelStore) Create(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
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

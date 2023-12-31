package db

import (
	"context"

	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStore interface {
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	GetHotels(context.Context, Map, *PaginationParams) ([]*types.Hotel, error)
	Create(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, Map, Map) error
	Delete(context.Context, string) error
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

func (store *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	var hotel *types.Hotel
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	dbhotel := store.coll.FindOne(ctx, bson.M{"_id": oid})
	if err := dbhotel.Decode(&hotel); err != nil {
		return nil, err
	}
	return hotel, nil
}

func (store *MongoHotelStore) GetHotels(ctx context.Context, filter Map, pagination *PaginationParams) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip((pagination.Page - 1) * pagination.Limit)
	opts.SetLimit(pagination.Limit)
	cur, err := store.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
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

func (store *MongoHotelStore) Update(ctx context.Context, filter, update Map) error {
	_, err := store.coll.UpdateOne(ctx, filter, update)
	return err
}

func (store *MongoHotelStore) Delete(ctx context.Context, id string) error {
	return nil
}

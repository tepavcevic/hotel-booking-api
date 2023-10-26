package db

import (
	"context"

	"github.com/tepavcevic/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	CreateRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(c *mongo.Client, dbname, collName string, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client: c,
		coll:   c.Database(dbname).Collection(collName),

		HotelStore: hotelStore,
	}
}

func (store *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	var rooms []*types.Room
	cur, err := store.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (store *MongoRoomStore) CreateRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := store.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := store.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}

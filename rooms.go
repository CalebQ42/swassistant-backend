package swassistantbackend

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetInvites(col *mongo.Collection, name string) (inv []RoomInv, err error) {
	cur, err := col.Find(context.TODO(), bson.D{{Key: "invites", Value: name}}, options.Find().SetProjection(bson.D{{Key: "name", Value: "1"}, {Key: "owner", Value: "1"}}))
	if err != nil {
		return
	}
	if cur.Err() == mongo.ErrNoDocuments {
		return
	}
	err = cur.All(context.TODO(), &inv)
	return
}

func GetRooms(col *mongo.Collection, name string) (rooms []Room, err error) {
	cur, err := col.Find(context.TODO(), bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "users", Value: name}}, bson.D{{Key: "owner", Value: name}}}}})
	if err != nil {
		return
	}
	if cur.Err() == mongo.ErrNoDocuments {
		return
	}
	err = cur.All(context.TODO(), &rooms)
	return
}

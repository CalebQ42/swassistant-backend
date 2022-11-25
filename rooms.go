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

func GetRooms(roomCol *mongo.Collection, edCol *mongo.Collection, name string) (rooms []Room, err error) {
	cur, err := roomCol.Find(context.TODO(), bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "users", Value: name}}, bson.D{{Key: "owner", Value: name}}}}})
	if err != nil {
		return
	}
	if cur.Err() == mongo.ErrNoDocuments {
		return
	}
	err = cur.All(context.TODO(), &rooms)
	for i := range rooms {
		if rooms[i].Owner != name {
			rooms[i].Invites = nil
			rooms[i].Declined = nil
		}
		cur, err = edCol.Find(context.TODO(), bson.D{{Key: "perm.name", Value: rooms[i].ID}}, options.Find().SetProjection(bson.D{{Key: "_id", Value: "1"}}))
		if err != nil {
			return
		}
		if cur.Err() == mongo.ErrNoDocuments {
			continue
		}
		err = cur.All(context.TODO(), &(rooms[i].Editables))
		if err != nil {
			return
		}
	}
	return
}

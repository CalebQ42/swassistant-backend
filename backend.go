package swassistantbackend

import (
	"context"
	"log"

	"github.com/CalebQ42/stupid-backend"
	"github.com/CalebQ42/stupid-backend/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SWBackend struct {
	stupid.App
	db *mongo.Database
}

func NewSWBackend(client *mongo.Client) *SWBackend {
	return &SWBackend{
		db: client.Database("swassistant"),
	}
}

func (s *SWBackend) Log() db.LogTable {
	return db.NewMongoTable(s.db.Collection("logs"))
}

func (s *SWBackend) Crashes() db.CrashTable {
	return db.NewMongoTable(s.db.Collection("crashes"))
}

func (s *SWBackend) IgnoreOldVersionCrashes() bool {
	return true
}

func (s *SWBackend) CurrentVersions() (out []string) {
	out = make([]string, 0)
	res, err := s.db.Collection("versions").Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("SWAssistant: Error getting versions:", err)
		return
	}
	var vers []struct {
		ID  string `bson:"_id"`
		Ver string `bson:"version"`
	}
	err = res.All(context.TODO(), &vers)
	if err != nil {
		log.Println("SWAssistant: Error decoding versions:", err)
		return
	}
	out = make([]string, len(vers))
	for i := range vers {
		out[i] = vers[i].Ver
	}
	return
}

func (s *SWBackend) Extention(req *stupid.Request) bool {
	if len(req.Path) == 0 {
		return true
	}
	switch req.Path[0] {
	case "rooms":
		return s.HandleRooms(req)
	case "profile":
		//TODO
	}
	return false
}

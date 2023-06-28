package swassistantbackend

import (
	"context"
	"log"
	"time"

	"github.com/CalebQ42/stupid-backend"
	"github.com/CalebQ42/stupid-backend/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SWBackend struct {
	db *mongo.Database
}

func NewSWBackend(client *mongo.Client) *SWBackend {
	go func() {
		for range time.Tick(time.Hour) {
			log.Println("SWAssistant: Deleting expired profiles")
			res, err := client.Database("swassistant").Collection("profiles").DeleteMany(context.TODO(), bson.M{"expiration": bson.M{"$lt": time.Now().Unix()}})
			if err == mongo.ErrNoDocuments {
				continue
			}
			log.Println("SWAssistant: Deleted", res.DeletedCount, "profiles")
		}
	}()
	return &SWBackend{
		db: client.Database("swassistant"),
	}
}

func (s *SWBackend) Logs() db.LogTable {
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

func (s *SWBackend) Extension(req *stupid.Request) bool {
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

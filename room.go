package swassistantbackend

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/CalebQ42/stupid-backend/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Room struct {
	ID       string   `json:"id" bson:"_id"`
	Name     string   `json:"name" bson:"name"`
	Owner    string   `json:"owner" bson:"owner"`
	Users    []string `json:"users" bson:"users"`
	Profiles []string `json:"profiles" bson:"profiles"`
}

func (s *SWBackend) HandleRooms(req *stupid.Request) bool {
	if len(req.Path) != 2 {
		req.Resp.WriteHeader(http.StatusBadRequest)
		return true
	}
	switch req.Path[1] {
	case "list":
		return s.ListRooms(req)
	case "new":
		return s.NewRoom(req)
	default:
		return s.GetRoom(req)
	}
}

func (s *SWBackend) ListRooms(req *stupid.Request) bool {
	if req.Method != http.MethodGet {
		req.Resp.WriteHeader(http.StatusBadRequest)
		return true
	} else if req.User == nil {
		req.Resp.WriteHeader(http.StatusUnauthorized)
		return true
	}
	out := make([]struct {
		ID    string `json:"id" bson:"_id"`
		Name  string `json:"name" bson:"name"`
		Owner string `json:"owner" bson:"owner"`
	}, 0)
	res, err := s.db.Collection("rooms").Find(context.TODO(), bson.M{"users": req.User.Username}, options.Find().SetProjection(bson.M{"_id": 1, "name": 1, "owner": 1}))
	if err != nil && err != mongo.ErrNoDocuments {
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	if err == nil {
		err = res.All(context.TODO(), &out)
		if err != nil {
			log.Println("SWAssistant: Error decoding room list:", err)
			req.Resp.WriteHeader(http.StatusInternalServerError)
			return true
		}
	}
	outDat, err := json.Marshal(out)
	if err != nil {
		log.Println("SWAssistant: Error encoding room list:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	_, err = req.Resp.Write(outDat)
	if err != nil {
		log.Println("SWAssistant: Error writing room list:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
	}
	return true
}

func (s *SWBackend) NewRoom(req *stupid.Request) bool {
	if req.Method != http.MethodPost || req.Query["name"] == nil || len(req.Query["name"]) != 1 || req.Query["name"][0] == "" {
		req.Resp.WriteHeader(http.StatusBadRequest)
		return true
	} else if req.User == nil {
		req.Resp.WriteHeader(http.StatusUnauthorized)
		return true
	}
	//TODO: check room name for unsavory words
	newRoom := Room{
		ID:       uuid.NewString(),
		Name:     req.Query["name"][0],
		Owner:    req.User.Username,
		Users:    []string{},
		Profiles: []string{},
	}
	_, err := s.db.Collection("rooms").InsertOne(context.TODO(), newRoom)
	if err != nil {
		log.Println("SWAssistant: Error creating room:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	out, err := json.Marshal(map[string]string{"id": newRoom.ID, "name": newRoom.Name})
	if err != nil {
		log.Println("SWAssistant: Error encoding new room:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	req.Resp.WriteHeader(http.StatusCreated)
	_, err = req.Resp.Write(out)
	if err != nil {
		log.Println("SWAssistant: Error writing new room:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
	}
	return true
}

func (s *SWBackend) GetRoom(req *stupid.Request) bool {
	if req.Method != http.MethodGet {
		req.Resp.WriteHeader(http.StatusBadRequest)
		return true
	} else if req.User == nil {
		req.Resp.WriteHeader(http.StatusUnauthorized)
		return true
	}
	res := s.db.Collection("rooms").FindOne(context.TODO(), bson.M{"_id": req.Path[1]})
	if res.Err() == mongo.ErrNoDocuments {
		req.Resp.WriteHeader(http.StatusNotFound)
		return true
	} else if res.Err() != nil {
		log.Println("SWAssistant: Error getting room:", res.Err())
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	r := Room{}
	err := res.Decode(&r)
	if err != nil {
		log.Println("SWAssistant: Error decoding room:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	out, err := json.Marshal(r)
	if err != nil {
		log.Println("SWAssistant: Error encoding room:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	_, err = req.Resp.Write(out)
	if err != nil {
		log.Println("SWAssistant: Error writing room:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
	}
	return true
}

package swassistantbackend

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/CalebQ42/stupid-backend"
	"github.com/lithammer/shortuuid/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *SWBackend) HandleProfiles(req *stupid.Request) bool {
	if len(req.Path) != 2 {
		req.Resp.WriteHeader(http.StatusBadRequest)
		return true
	}
	switch req.Path[1] {
	case "upload":
		return s.UploadProfile(req)
	default:
		return s.GetProfile(req)
	}
}

type UploadedProf struct {
	Profile    map[string]any `json:"profile" bson:"profile"`
	ID         string         `json:"id" bson:"_id"`
	Type       string         `json:"type" bson:"type"`
	Expiration int64          `json:"expiration" bson:"expiration"`
}

func (s *SWBackend) UploadProfile(req *stupid.Request) bool {
	if req.Method != http.MethodPost || req.Query["type"] == nil || len(req.Query["type"]) != 1 {
		req.Resp.WriteHeader(http.StatusBadRequest)
		return true
	}
	profType := req.Query["type"][0]
	if profType != "character" && profType != "vehicle" && profType != "minion" {
		req.Resp.WriteHeader(http.StatusBadRequest)
		return true
	}
	if req.Body == nil {
		req.Resp.WriteHeader(http.StatusBadRequest)
		return true
	}
	data, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		log.Println("SWAssistant: Error reading incoming profile:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	} else if len(data) == 0 {
		req.Resp.WriteHeader(http.StatusBadRequest)
		return true
	} else if len(data) > 5242880 { // 5MB
		req.Resp.WriteHeader(http.StatusRequestEntityTooLarge)
		return true
	}
	prof := make(map[string]any)
	err = json.Unmarshal(data, &prof)
	if err != nil {
		log.Println("SWAssistant: Error decoding incoming profile:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	toUpload := UploadedProf{
		ID:         shortuuid.New(),
		Expiration: time.Now().Add(time.Hour * 12).Round(time.Hour).Unix(),
		Type:       profType,
		Profile:    prof,
	}
	_, err = s.db.Collection("profiles").InsertOne(context.TODO(), toUpload)
	if err != nil {
		log.Println("SWAssistant: Error inserting profile:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	out, err := json.Marshal(map[string]any{"id": toUpload.ID, "expiration": toUpload.Expiration})
	if err != nil {
		log.Println("SWAssistant: Error encoding profile response:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	req.Resp.WriteHeader(http.StatusCreated)
	_, err = req.Resp.Write(out)
	if err != nil {
		log.Println("SWAssistant: Error writing profile response:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	return true
}

func (s *SWBackend) GetProfile(req *stupid.Request) bool {
	if req.Method != http.MethodGet {
		req.Resp.WriteHeader(http.StatusBadRequest)
		return true
	}
	res := s.db.Collection("profiles").FindOne(context.TODO(), bson.M{"_id": req.Path[1]})
	if res.Err() == mongo.ErrNoDocuments {
		req.Resp.WriteHeader(http.StatusNotFound)
		return true
	} else if res.Err() != nil {
		log.Println("SWAssistant: Error getting profile:", res.Err())
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	var prof UploadedProf
	err := res.Decode(&prof)
	if err != nil {
		log.Println("SWAssistant: Error decoding profile:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	out, err := json.Marshal(prof.Profile)
	if err != nil {
		log.Println("SWAssistant: Error encoding profile:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
		return true
	}
	_, err = req.Resp.Write(out)
	if err != nil {
		log.Println("SWAssistant: Error writing profile:", err)
		req.Resp.WriteHeader(http.StatusInternalServerError)
	}
	return true
}

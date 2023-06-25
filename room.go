package swassistantbackend

import "github.com/CalebQ42/stupid-backend"

type Room struct {
	ID       string   `json:"_id" bson:"_id"`
	Name     string   `json:"name" bson:"name"`
	Owner    string   `json:"owner" bson:"owner"`
	Users    []string `json:"users" bson:"users"`
	Profiles []string `json:"profiles" bson:"profiles"`
}

func (s *SWBackend) HandleRooms(req *stupid.Request) bool {
	return true
}

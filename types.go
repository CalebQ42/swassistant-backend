package swassistantbackend

type EditablePerm struct {
	Name  string
	Room  bool
	Write bool
}

type RoomInv struct {
	ID    string `json:"_id" bson:"_id"`
	Name  string
	Owner string
}

type Room struct {
	ID        string `json:"_id" bson:"_id"`
	Name      string
	Owner     string
	Users     []string
	Invites   []string
	Declined  []string
	Editables []string `json:"-" bson:"-"`
}

func ValidEditable(m map[string]any) (ok bool) {
	ok = true
	if _, ok = m["_id"]; !ok {
		return
	}
	if _, ok = m["owner"]; !ok {
		return
	}
	if _, ok = m["perm"]; !ok {
		return
	}
	return
}

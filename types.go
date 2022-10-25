package swassistant

type room struct {
	ID string `json:"_id" bson:"_id"`
    Name string `json:"name" bson:"name"`
    Owner string `json:"owner" bson:"owner"`
    Users []string `json:"users" bson:"users"`
    Editables []string `json:"editables" bson:"editables"`
}

package persistence

import "gopkg.in/mgo.v2/bson"

type DataBaseHandler interface {
	AddEvent(Event) ([]byte, error)
	FindEvent([]byte) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllEventAvialable() ([]Event, error)
	AddLocation(Event) ([]byte, error)
}

type Event struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Duration  int
	StartDate int64
	EndDate   int64
	Location  Location
}
type Location struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Address   string
	Country   string
	OpenTime  string
	CloseTime string
	Hal       []Hall
}
type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}

type User struct {
	ID        bson.ObjectId `json:"bson_id"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	Password  string        `json:"password"`
	CreatedAt string        `json:"created_at"`
}

type UserRequest struct {
	Username string `json:"username" binding:"required, alphanum"`
	Password string `json:"password" binding:"required, min=7"`
}

type UserResponse struct {
	ID        bson.ObjectId `json:"bson_id"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	Password  string        `json:"password"`
	CreatedAt string        `json:"created_at"`
}
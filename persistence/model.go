package persistence

import "gopkg.in/mgo.v2/bson"

type DataBaseHandler interface{
	AddEvent(Event) ([]byte, error)
	FindEvent([]byte) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllEventAvialable() ([]Event, error)
}

type Event struct{
	ID  bson.ObjectId `bson:"_id"`
	Name string
	Duration int
    StartDate int64
	EndDate int64 
	Location  Location
}
type Location struct{
	ID  bson.ObjectId `bson:"_id"`
	Name string
	Address string
	Country string
	OpenTime string
	CloseTime string
	Hal []Hall
}
type Hall struct{
	Name string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int `json:"capacity"`

}
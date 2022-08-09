package database

import (
	"myevent/persistence"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)



const(
	DB = "myevent"
    USERS = "users"
	EVENTS = "events"
)
type MongoLayer struct{
	mongoLayer *mgo.Session
}

func NewMongoLayer(connection string) (*MongoLayer, error){
 s, err := mgo.Dial(connection)
 if err != nil{
	return nil,err
 }
 return &MongoLayer{
	mongoLayer: s,
 }, err
}

func (monoLayer *MongoLayer) getFresh() *mgo.Session {
 return monoLayer.mongoLayer.Copy()
}

func (monolayer *MongoLayer) AddEvent(e persistence.Event) ([]byte, error) {
	s := monolayer.getFresh()
	defer s.Close()
	if !e.ID.Valid(){
         e.ID = bson.NewObjectId()
	}
	if !e.Location.ID.Valid(){
		e.Location.ID = bson.NewObjectId()
	}
	return []byte(e.ID), s.DB(DB).C(EVENTS).Insert(e)
} 

func (mongoLayer *MongoLayer) FindEvent(id []byte) (persistence.Event, error) {
  s := mongoLayer.getFresh()
  defer s.Close()
  e := persistence.Event{}
  err := s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e)
  return e,err
}

func (monolayer *MongoLayer) FindEventByName(name string) (persistence.Event, error) {
	s := monolayer.getFresh()
	defer s.Close()
	e := persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name":name}).One(&e)
	return e, err
	
}
func (monolayer *MongoLayer) FindAllEventAvialable() ([]persistence.Event, error) {
	s := monolayer.getFresh()
	defer s.Close()
	events := []persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(nil).All(&events)
	return events, err
}


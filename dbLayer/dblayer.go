package dblayer

import (
	"myevent/database"
	"myevent/persistence"
)


type DATATYPE string


const (
	MONGODB DATATYPE = "mongodb"
	DYNAMO DATATYPE = "dynamodb"
)


func NewPersistenceLayer(options DATATYPE, connection string) (persistence.DataBaseHandler, error){
  switch options{
  case MONGODB:
	return database.NewMongoLayer(connection)
  }
  return nil, nil
}

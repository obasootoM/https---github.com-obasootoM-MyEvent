package configuration

import (
	"encoding/json"
	"fmt"
	dblayer "myevent/dbLayer"
	"os"
)

var (
	DATATYPEDEFAULT       = dblayer.MONGODB
	DATACONNECTIONDEFAULT = "mongodb:127.0.0.1"
	RESTFULDEFAULT        = "localhost:8080"
)

type ServiceConfig struct {
	Databasetype      dblayer.DATATYPE `json:"databasetype"`
	DatabaseConnection string           `json:"databaseconnection"`
	RestfulEndpoint   string           `json:"restfulendpoint"`
}

func NewServiceConfig(fileName string) (*ServiceConfig, error) {
	config := &ServiceConfig{
		DATATYPEDEFAULT,
		DATACONNECTIONDEFAULT,
		RESTFULDEFAULT,
	}
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("file cannot be open, need to continue")
		return config, err
	}
	decode := json.NewDecoder(file)
	err = decode.Decode(&config)
	return config, err
}

package configuration

import (
	"encoding/json"
	"fmt"
	dblayer "myevent/dbLayer"
	"os"
)

var (
	DATATYPEDEFAULT       = dblayer.MONGODB
	DATACONNECTIONDEFAULT = "mongodb://127.0.0.1"
	RESTFULDEFAULT        = "localhost:9191"
	RESTFULDEFAULTLS      = "localhost:8181"
	AMQPMESSAGEBROKER     = "amqp://guest:guest@localhost:8282"
)

type ServiceConfig struct {
	Databasetype       dblayer.DATATYPE `json:"databasetype"`
	DatabaseConnection string           `json:"databaseconnection"`
	RestfulEndpoint    string           `json:"restfulendpoint"`
	RestfulEndpointTls string           `json:"restfulEndpointTls"`
	AmqpMessageBroker  string           `json:"amqpmessagebroker"`
}

func NewServiceConfig(fileName string) (*ServiceConfig, error) {
	config := &ServiceConfig{
		DATATYPEDEFAULT,
		DATACONNECTIONDEFAULT,
		RESTFULDEFAULT,
		RESTFULDEFAULTLS,
		AMQPMESSAGEBROKER,
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("configuration file not found, need to continue")
		return config, err
	}
	decode := json.NewDecoder(file)
	err = decode.Decode(&config)
	if broker := os.Getenv("AMQPMESSAGEBROKER"); broker != "" {
		config.AmqpMessageBroker = broker
	}
	if err != nil {
		fmt.Println("cannot decode config")
	}
	return config, err
}

package main

import (
	"flag"
	"fmt"
	"myevent/api"
	"myevent/bookingservice/listner"
	"myevent/configuration"
	dblayer "myevent/dbLayer"

	amqp_test "myevent/lib/mesqp/amqp"

	"github.com/streadway/amqp"
)

func main() {
	confPath := flag.String("conf", `./configuration/config.json`, "set the path to configuration json file")
	flag.Parse()
	config, _ := configuration.NewServiceConfig(*confPath)
	fmt.Println("connecting to database")
	dbHandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DatabaseConnection)
	connection, err := amqp.Dial(config.AmqpMessageBroker)
	if err != nil {
		panic("could not establish amqp connection" + err.Error())
	}
	defer connection.Close()
	eventListener, err := amqp_test.NewAmpqEventListner(connection, "")
	if err != nil {
		panic(err)
	}
	emitter, err := amqp_test.NewAmpqEventEmitter(connection)
	if err != nil {
		panic(err)
	}
	process := &listner.EventProcessor{
		EventListner: eventListener,
		Database:     dbHandler,
	}
	go process.ProcessEvent()
	api.ServiceApi(config.RestfulEndpoint, config.RestfulEndpointTls, dbHandler, emitter)
}
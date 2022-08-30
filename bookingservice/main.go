package main

import (
	"flag"
	"fmt"
	"log"
	"myevent/api"
	"myevent/configuration"
	dblayer "myevent/dbLayer"
	amqp_test "myevent/lib/amqp"
	listner_test "myevent/lib/listner"

	"github.com/streadway/amqp"
	"myevent/bookingservice/listner"
)



func main() {
	path := flag.String("conf", `.\configuration\config.json`, "flag to set  path to configuration json file")
	flag.Parse()
	configu, _ := configuration.NewServiceConfig(*path)
	connection, err := amqp.Dial(configu.AmqpMessageBroker)
	if err != nil {
		log.Fatal("cannot secure connection to rabbitmq" + err.Error())
	}
	defer connection.Close()
	fmt.Println("connecting to database")
    dbHandler, _ := dblayer.NewPersistenceLayer(configu.Databasetype, configu.DatabaseConnection)
	list, err := listner_test.NewAmpqEventListner(connection, "")
	if err != nil{
      log.Fatal("cannot create listner" + err.Error())
	}
	emitter, err := amqp_test.NewAmpqEventEmitter(connection)
	if err != nil {
		log.Fatal("")
	}

	processor := listner.EventProcessor{
		EventListner: list,
		Database: dbHandler,
	}
	go processor.ProcessEvent()
	
	api.ServiceApi(configu.RestfulEndpoint,configu.RestfulEndpointTls, dbHandler, emitter)
}
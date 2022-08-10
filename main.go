package main

import (
	"flag"
	"fmt"
	"log"
	"myevent/api"
	"myevent/configuration"
	dblayer "myevent/dbLayer"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	amqpURL := os.Getenv("amqp")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672"
	}
	confPath := flag.String("conf", `./configuration/config.json`, "set the path to configuration json file")
	flag.Parse()
	config, _ := configuration.NewServiceConfig(*confPath)
	fmt.Println("connecting to database")
	dbHandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DatabaseConnection)
	connection, err := amqp.Dial(amqpURL)
	if err != nil {
		panic("could not establish amqp connection" + err.Error())
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic("cannot connect to channel" + err.Error())
	}
	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	message := amqp.Publishing{
		Body: []byte("Hello World"),
	}
	err = channel.Publish("event", "some-routing-key", false, false, message)
	if err != nil {
		panic("cannot publish channel" + err.Error())
	}
	_, err = channel.QueueDeclare("my_queue",true,false,false,false,nil)
	if err != nil {
		panic("error while declaring queue" + err.Error())
	}
	 err = channel.QueueBind("event","","",true,nil)
	 if err != nil {
		panic("error when declaring bind to queue" + err.Error())
	 }
	 msg, err := channel.Consume("event","",false,false,false,false,nil)
	 if err != nil{
		panic("error cannot consume" + err.Error())
	 }
    for message := range msg{
        fmt.Println("message recieved" + string(message.Body))
		message.Ack(false)
	}
	httpChanServe, httpChanServeTls := api.ServiceApi(config.RestfulEndpoint, config.RestfulEndpointTls, dbHandler)
	select {
	case err := <-httpChanServe:
		log.Fatal("HTTP ERROR", err)
	case err := <-httpChanServeTls:
		log.Fatal("HTTPS", err)
	}

}

package main

import (
	"flag"
	"fmt"
	"log"
	"myevent/api"
	"myevent/configuration"
	dblayer "myevent/dbLayer"
	
     "myevent/lib/amqp"
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
	emitter, err := amqp_test.NewAmpqEventEmitter(connection)
	if err != nil{
       panic(err)
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
		Body: []byte("\nHello World"),
	}
	err = channel.Publish("events", "some-routing-key", false, false, message)
	if err != nil {
		panic("cannot publish channel" + err.Error())
	}
	_, err = channel.QueueDeclare("my_queue",true,false,false,false,nil)
	if err != nil {
		panic("error while declaring queue" + err.Error())
	}
	 err = channel.QueueBind("my_queue","#","events",false,nil)
	 if err != nil {
		panic("error when declaring bind to queue" + err.Error())
	 }
	 msg, err := channel.Consume("my_queue","",false,false,false,false,nil)
	 if err != nil{
		panic("error cannot consume" + err.Error())
	 }
    for message := range msg{
        fmt.Println("message recieved" + string(message.Body))
		message.Ack(false)
	}
	httpChanServe, httpChanServeTls := api.ServiceApi(config.RestfulEndpoint, config.RestfulEndpointTls, dbHandler,emitter)
	select {
	case err := <-httpChanServe:
		log.Fatal("HTTP ERROR", err)
	case err := <-httpChanServeTls:
		log.Fatal("HTTPS", err)
	}

}

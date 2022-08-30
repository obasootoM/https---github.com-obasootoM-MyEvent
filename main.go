package main

import (
	"flag"
	"fmt"
	"log"
	"myevent/api"
	"myevent/configuration"
	dblayer "myevent/dbLayer"
	amqp_test "myevent/lib/amqp"

	"github.com/streadway/amqp"
)

func main() {
	path := flag.String("conf", `.\configuration\config.json`, "flag to set  path to configuration json file")
	flag.Parse()
	configu, _ := configuration.NewServiceConfig(*path)
	connection, err := amqp.Dial("amqp://guest:guest@localhost:8282/")
	if err != nil {
		log.Fatal("cannot secure connection to rabbitmq" + err.Error())
	}
	defer connection.Close()
	fmt.Println("connecting to database")
	channel, err := connection.Channel()
	if err != nil {
		log.Fatal("connection err" + err.Error())
	}
	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal("cannot exchange declare " + err.Error())
	}
	p, err := channel.QueueDeclare("my_event", true, false, false, false, nil)
	if err != nil {
		log.Fatal("cannot declare queue" + err.Error())
	}
	err = channel.QueueBind(p.Name, "#", "events",false,nil)
	if err != nil {
		log.Fatal("cannot bind to queue" + err.Error())
	}
	body := "hello world"
	msg := amqp.Publishing{
		Body:        []byte(body),
		ContentType: "text/body",
	}
	err = channel.Publish("events", "key", false, false, msg)
	if err != nil {
		log.Fatal("cannot publish " + err.Error())
	}
	msgs, err := channel.Consume(p.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("cannot consume " + err.Error())
	}
	for mesg := range msgs {
		fmt.Println("message recieved" + string(mesg.Body))
		mesg.Ack(false)
	}

	dbHandler, _ := dblayer.NewPersistenceLayer(configu.Databasetype, configu.DatabaseConnection)
	emiter, err := amqp_test.NewAmpqEventEmitter(connection)
	if err != nil {
		log.Fatal("cannot create new emitter" + err.Error())
	}

	httpServer, httpServerTls := api.ServiceApi(configu.RestfulEndpoint, configu.RestfulEndpointTls, dbHandler, emiter)
	select {
	case err := <-httpServer:
		log.Fatal("http err" + err.Error())
	case err := <-httpServerTls:
		log.Fatal("https err" + err.Error())
	}
}

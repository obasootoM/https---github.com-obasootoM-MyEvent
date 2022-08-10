package main

import (
	"flag"
	"fmt"
	"log"
	"myevent/api"
	"myevent/configuration"
	dblayer "myevent/dbLayer"
)


func main() {
	confPath := flag.String("conf",`./configuration/config.json`,"set the path to configuration json file")
	flag.Parse()
	config, _ := configuration.NewServiceConfig(*confPath)
	fmt.Println("connecting to database")
	dbHandler, _ := dblayer.NewPersistenceLayer(config.Databasetype,config.DatabaseConnection)
	log.Fatal(api.ServiceApi(config.RestfulEndpoint,dbHandler),)

}
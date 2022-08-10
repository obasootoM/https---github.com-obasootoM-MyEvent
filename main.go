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

	httpChanServe,httpChanServeTls :=api.ServiceApi(config.RestfulEndpoint,config.RestfulEndpointTls,dbHandler)
	select {
	case err := <- httpChanServe:
		log.Fatal("HTTP ERROR", err)
	case err := <- httpChanServeTls:
		log.Fatal("HTTPS",err)
	}

}
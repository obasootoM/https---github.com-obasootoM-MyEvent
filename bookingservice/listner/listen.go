package listner

import (
	"log"
	"myevent/contract"
	"myevent/lib/mesqp"
	"myevent/persistence"

	"gopkg.in/mgo.v2/bson"
)

type EventProcessor struct {
	EventListner mesqp.EventListner
	Database     persistence.DataBaseHandler
}

func (e *EventProcessor) ProcessEvent() error {
	log.Fatal("listening to event...")
    recieved, error, err := e.EventListner.Listen("event.created")
	if err != nil {
		return err
	}
	for {
		select {
		case evt := <- recieved:
			e.handleEvent(evt)
		case err := <- error :
			log.Printf("received error when processing msg %s",err)	
		}
	
	}
}

func (p *EventProcessor) handleEvent(evt mesqp.Event)  {
    switch e := evt.(type) {
	case *contract.EventCreatedEvent:
		log.Printf("event %s created %s", e.Id,e)
		p.Database.AddEvent(persistence.Event{ID: bson.ObjectId(e.Id)})
	case *contract.LocationCreatedEvent:
		log.Printf("event %s created %s", e.Id,e)
		p.Database.AddLocation(persistence.Event{})
	}
}

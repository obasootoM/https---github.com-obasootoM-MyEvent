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

func (p *EventProcessor) ProcessEvent() error {
	log.Println("listening to event...")
	recieved, errors, err := p.EventListner.Listen("event.created")
	if err != nil {
		return err
	}
	for {
		select {
		case evt := <-recieved:
			p.handleEvent(evt)
		case err := <-errors:
			log.Printf("received error when processing msg %s", err)
		}

	}
}

func (p *EventProcessor) handleEvent(event mesqp.Event) {
	switch e := event.(type) {
	case *contract.EventCreatedEvent:
		log.Printf("event %s created %s", e.Id, e)
		p.Database.AddEvent(persistence.Event{ID: bson.ObjectId(e.Id)})
	case *contract.LocationCreatedEvent:
		log.Printf("event %s created %s", e.Id, e)
		p.Database.AddLocation(persistence.Event{ID: bson.ObjectId(e.Id)})
	}
}


package amqp_test

import (
	"encoding/json"


	"fmt"
	"myevent/contract"
	"myevent/lib/mesqp"

	"github.com/streadway/amqp"
)



type ampqEventListner struct {
	connection *amqp.Connection
	query      string
}

func (a *ampqEventListner) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	_, err = channel.QueueDeclare(a.query, true, false, false, false, nil)
	return err
}
func NewAmpqEventListner(con *amqp.Connection, query string) (mesqp.EventListner, error) {
	listener := &ampqEventListner{
		connection: con,
		query:      query,
	}
	err := listener.setup()
	if err != nil {
		return nil, err
	}
	return listener, err
}

func (a *ampqEventListner) Listen(eventNames ...string) (<-chan mesqp.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	defer channel.Close()
	for _, eventName := range eventNames {
		if err := channel.QueueBind(a.query, eventName, "events", false, nil); err != nil {
			return nil, nil, err
		}
	}
	msgs, err := channel.Consume(a.query, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}
	events := make(chan mesqp.Event)
	errors := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errors <- fmt.Errorf("msg did not contain x-event-name header")
				msg.Nack(true, false)
				continue
			}
			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("x-event-name header is not a string %t", rawEventName)
				msg.Nack(true, false)
				continue
			}
			var event mesqp.Event
			switch eventName {
			case "event.created":
				event = new(contract.EventCreatedEvent)
            default:
				errors <- fmt.Errorf("event type %s is unknown",eventName)
				continue
			}
			err := json.Unmarshal(msg.Body,event)
			if err != nil {
				errors <- err
				continue
			}
			events <- event
		}
	}()
	return events, errors, nil
}

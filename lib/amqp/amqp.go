package amqp_test

import (
	"encoding/json"
	"fmt"
	"myevent/contract"
	"myevent/lib/mesqp"

	"github.com/streadway/amqp"
)

type amqpEventEmitter struct {
	connection *amqp.Connection
}
type ampqEventListner struct {
	connection *amqp.Connection
	query      string
}

// Listen implements mesqp.EventListner
func (a *ampqEventListner) Listen(eventNames ...string) (<-chan mesqp.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	defer channel.Close()
	for _, eventName := range eventNames {
		if err := channel.QueueBind(a.query, eventName, "event", false, nil); err != nil {
			return nil, nil, err
		}
	}
	msgs, err := channel.Consume(a.query, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}
	events := make(chan mesqp.Event)
	errorss := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errorss <- fmt.Errorf("msg did not contain x-event-header")
				msg.Nack(false, false)
				continue
			}
			eventName, ok := rawEventName.(string)
			if !ok {
				errorss <- fmt.Errorf("x-event-header is not a string %t", rawEventName)
				msg.Nack(false, false)
				continue
			}
			var event mesqp.Event
			switch eventName {
			case "event.created":
				event = new(contract.EventCreatedEvent)
            default:
				errorss <- fmt.Errorf("event type %s is unknown",eventName)
				continue
			}
			err := json.Unmarshal(msg.Body,event)
			if err != nil {
				errorss <- err
				continue
			}
			events <- event
		}
	}()
	return events, errorss, nil
}

func (a *amqpEventEmitter) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return channel.ExchangeDeclare("event", "topic", true, false, false, false, nil)
}
func (a *ampqEventListner) set() error {
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
	err := listener.set()
	if err != nil {
		return nil, err
	}
	return listener, err
}
func NewAmpqEventEmitter(conn *amqp.Connection) (mesqp.EventEmmiter, error) {
	emiter := &amqpEventEmitter{
		connection: conn,
	}
	err := emiter.setup()
	if err != nil {
		return nil, err
	}
	return emiter, nil
}

func (a *amqpEventEmitter) Emit(event mesqp.Event) error {
	decodeJ, err := json.Marshal(event)
	if err != nil {
		return err
	}
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": event.EventName()},
		Body:        decodeJ,
		ContentType: "application/json",
	}
	return channel.Publish(
		"event",
		event.EventName(),
		false,
		false,
		msg,
	)
}

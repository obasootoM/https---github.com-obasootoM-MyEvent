package amqp_test

import (
	"encoding/json"
	"myevent/lib/mesqp"

	"github.com/streadway/amqp"
)




type amqpEventEmitter struct{
	connection *amqp.Connection
}

func (a *amqpEventEmitter) setup() error {
   channel, err := a.connection.Channel()
   if err != nil {
	return err
   }
   defer channel.Close()
   return channel.ExchangeDeclare("event","topic", true,false,false,false, nil)
}

func NewAmpqEventEmitter(conn *amqp.Connection) (mesqp.EventEmmiter, error) {
	emiter := &amqpEventEmitter{
      connection: conn,
	}
	err := emiter.setup()
	if err != nil {
     return nil,err
	}
	return emiter,nil
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
	Headers: amqp.Table{"x-event-name":event.EventName()},
	Body: decodeJ,
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
package mesqp


type Event interface{
   EventName() string
}
type EventEmmiter interface{
	Emit(event Event) error
}

type EventListner interface{
	Listen(eventNames ...string) (<-chan Event, <-chan error,error)
}
package mesqp


type Event interface{
   EventName() string
}
type EventEmmiter interface{
	Emit(event Event) error
}
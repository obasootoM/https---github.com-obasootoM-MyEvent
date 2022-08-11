package contract

import "time"


type EventCreatedEvent struct{
	Id string `json:"id"`
	Name string `json:"name"`
	LocationId string `json:"location_id"`
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
}

func (e *EventCreatedEvent) EventName() string{
   return "event_created"
}
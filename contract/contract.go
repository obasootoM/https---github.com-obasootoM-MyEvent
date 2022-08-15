package contract

import "time"

type EventCreatedEvent struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	LocationId string    `json:"location_id"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
}
type LocationCreatedEvent struct {
	Id      string    `json:"id"`
	Name    string    `josn:"name"`
	Adress  string    `json:"address"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Country string    `json:"country"`
}

func (e *EventCreatedEvent) EventName() string {
	return "event_created"
}
func (l *LocationCreatedEvent) EventName() string {
	return "event_location"
}

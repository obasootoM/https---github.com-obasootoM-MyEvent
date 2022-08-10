package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"myevent/persistence"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Service struct {
	DbHandler persistence.DataBaseHandler
}

func (eh *Service) findEventHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{error: cannot search describe events,you can either search by id or name}`)
		return
	}
	searchKey, ok := vars["search"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{error: no search key found, you can either search by id or name}`)
		return
	}
	var persist persistence.Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		persist, err = eh.DbHandler.FindEventByName(searchKey)
	case "id":
		id, err := hex.DecodeString(searchKey)
		if err == nil {
			persist,err = eh.DbHandler.FindEvent(id)

		}
	}
	if err != nil {
		fmt.Fprintf(w, "error %s", err)
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF8")
	err = json.NewEncoder(w).Encode(&persist)
	if err  != nil{
		fmt.Fprintf(w,"{error: cannot read from json%s}",err)
		return
	}

}
func (eh *Service) allEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.DbHandler.FindAllEventAvialable()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occur when trying to find all the events %s", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occur when trying to encode to json %s",err)
		return
	}
}
func (eh *Service) newEventHandler(w http.ResponseWriter, r *http.Request) {
	events := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occur when trying to decode from json%s", err)
		return
	}
	id, err := eh.DbHandler.AddEvent(events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occur while persisting event%d %s", id, err)
		return
	}
}

func NewService(DbHandler persistence.DataBaseHandler) *Service {
	return &Service{
		DbHandler: DbHandler,
	}
}

func ServiceApi(endpoint,tlsEndpoint string, dbHandler persistence.DataBaseHandler)(chan error,chan error) {
	r := mux.NewRouter()
	handler := NewService(dbHandler)
	eventRouter := r.PathPrefix("/events").Subrouter()
	eventRouter.Methods("Get").Path("").HandlerFunc(handler.findEventHandler)
	eventRouter.Methods("Get").Path("").HandlerFunc(handler.allEventHandler)
	eventRouter.Methods("Post").Path("").HandlerFunc(handler.newEventHandler)
	httpChanServe := make(chan error)
	httpChanServeTls := make(chan error)
	go func() {
      httpChanServeTls <- http.ListenAndServeTLS(tlsEndpoint,"cert.pem","key.pem",r)
	}()
	go func() {
		httpChanServe <- http.ListenAndServe(endpoint, r)
	}()
	return httpChanServe,httpChanServeTls
		
}

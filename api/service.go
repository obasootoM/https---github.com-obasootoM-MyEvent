package api

import (
	"myevent/persistence"
	"net/http"

	"github.com/gorilla/mux"
)


type Service struct{
	DbHandler persistence.DataBaseHandler
}

func (eh *Service) findEventHandler(w http.ResponseWriter, r http.Request){}
func (eh *Service) allEventHandler(w http.ResponseWriter, r http.Request) {}
func (eh *Service) newEventHandler(w http.ResponseWriter, r http.Request) {}

func  NewService(DbHandler persistence.DataBaseHandler)*Service{
	
	handler := &Service{}
	r := mux.NewRouter()
	eventRouter := r.PathPrefix("/").Subrouter()
	eventRouter.Methods("Get").Path("").HandlerFunc(handler.allEventHandler())
    eventRouter.Methods("Post").Path("").HandlerFunc(handler.findEventHandler()) 
	eventRouter.Methods("Get").Path("").HandlerFunc(handler.newEventHandler())
	return &Service{
		DbHandler: DbHandler,
	}
}
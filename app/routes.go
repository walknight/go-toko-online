package app

import (
	"github.com/gorilla/mux"
	"github.com/walknight/gotoko/app/controllers"
)

func (server *Server) InitializeRoutes() {
	//set router
	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}

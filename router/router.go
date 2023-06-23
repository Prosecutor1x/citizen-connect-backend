package router

import (
	"github.com/Prosecutor1x/citizen-connect-frontend/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/addData", controllers.AddDataHandler).Methods("POST")

	return router
}

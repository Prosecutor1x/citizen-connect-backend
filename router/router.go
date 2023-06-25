package router

import (
	"github.com/Prosecutor1x/citizen-connect-frontend/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/addNewIssue", controllers.AddDataHandler).Methods("POST")
	router.HandleFunc("/api/deleteIssue/{id}", controllers.DeleteDataHandler).Methods("DELETE")
	router.HandleFunc("/api/fetchIssue", controllers.FetchAllDataHandler).Methods("GET")
	router.HandleFunc("/api/fetchIssue/{id}", controllers.FetchSingleDataHandler).Methods("GET")
	router.HandleFunc("/api/updateIssue/{id}", controllers.UpdateDataHandler).Methods("PUT")
	router.HandleFunc("/api/sendOtp", controllers.SendOtp).Methods("POST")
	router.HandleFunc("/api/verifyOtp", controllers.VerifyOtp).Methods("POST")

	return router
}

// new issue list
// issue delete
// issue fetch
//single issue fetch
// issue update

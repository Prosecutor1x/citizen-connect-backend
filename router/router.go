package router

import (
	authhandler "github.com/Prosecutor1x/citizen-connect-frontend/handlers/auth_handler"
	"github.com/Prosecutor1x/citizen-connect-frontend/handlers/issue_handler"
	"github.com/Prosecutor1x/citizen-connect-frontend/handlers/user_handler"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// issue routes
	router.HandleFunc("/api/addNewIssue", issue_handler.CreateIssueHandler).Methods("POST")
	router.HandleFunc("/api/fetchIssue", issue_handler.FetchAllIssueHandler).Methods("GET")
	router.HandleFunc("/api/deleteIssue/{id}", issue_handler.DeleteIssueHandler).Methods("DELETE")
	router.HandleFunc("/api/fetchIssue/{id}", issue_handler.FetchSingleIssueHandler).Methods("GET")
	router.HandleFunc("/api/updateIssue/{id}", issue_handler.UpdateIssueHandler).Methods("PUT")

	//auth routes
	router.HandleFunc("/api/sendOtp", authhandler.SendOtp).Methods("POST")
	router.HandleFunc("/api/verifyOtp", authhandler.VerifyOtp).Methods("POST")

	// user routes
	router.HandleFunc("/api/createUser", user_handler.CreateUserHandler).Methods("POST")
	router.HandleFunc("/api/checkUser", user_handler.CheckUserExist).Methods("POST")

	return router
}

// new issue list
// issue delete
// issue fetch
//single issue fetch
// issue update

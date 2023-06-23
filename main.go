package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Prosecutor1x/citizen-connect-frontend/router"
)

func main() {
	fmt.Println("API for citizen-connect app")
	fmt.Println("Server stating in port 4000")
	r := router.Router()
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Server started in port 4000")
}

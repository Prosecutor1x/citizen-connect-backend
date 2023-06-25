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

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// 	os.Exit(1)
	// }

	// accountSid = os.Getenv("TWILIO_ACCOUNT_SID")
	// authToken = os.Getenv("TWILIO_AUTH_TOKEN")

	// client = twilio.NewRestClientWithParams(twilio.RestClientParams{
	// 	Username: accountSid,
	// 	Password: authToken,
	// })

	r := router.Router()
	log.Fatal(http.ListenAndServe(":4001", r))
	fmt.Println("Server started in port 4000")
}

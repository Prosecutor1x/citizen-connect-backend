package authhandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

func SendOtp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := godotenv.Load(".env.local")
	if err != nil {
		// send error to client
		response := map[string]interface{}{
			"message": "Unable to send otp",
		}

		// Convert the response to JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			// Handle JSON marshaling error
			log.Fatal("Error marshaling JSON response: ", err)
		}

		// Send the response
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		w.Write(jsonResponse)                         // Write the JSON response to the client
		log.Fatal("Error loading .env.local file")
	}

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	var requestBody struct {
		Phone string `json:"phone"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		response := map[string]interface{}{
			"message": "Unable to send otp",
		}
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		json.NewEncoder(w).Encode(response)

		log.Fatal("Error decoding request body: ", err)
	}

	phone := requestBody.Phone
	if phone != "" {
		resp, serResp := generateAndSendOtp(accountSid, authToken, phone)
		if resp == "Error" {
			response := map[string]interface{}{
				"message": "Unable to send otp",
			}
			w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
			json.NewEncoder(w).Encode(response)
		} else {
			response := map[string]interface{}{
				"message":              "Otp sent successfully",
				"verificationResponse": resp,
				"serviceResponseParam": serResp,
			}
			w.WriteHeader(http.StatusOK) // Set the HTTP status code
			json.NewEncoder(w).Encode(response)
		}
	} else {
		response := map[string]interface{}{
			"message": "Unable to send otp",
		}
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		json.NewEncoder(w).Encode(response)
	}

}

func generateAndSendOtp(accountSid string, authToken string, phone string) (string, string) {

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   accountSid,
		Password:   authToken,
		AccountSid: accountSid,
	})

	_CreateServiceParams := &verify.CreateServiceParams{}
	_CreateServiceParams.SetFriendlyName("Bongobasi")

	CreateServiceResp, err := client.VerifyV2.CreateService(_CreateServiceParams)
	if err != nil {
		return "Error", "Error"
	}

	_CreateVerificationParams := &verify.CreateVerificationParams{}
	_CreateVerificationParams.SetTo(phone)
	_CreateVerificationParams.SetChannel("sms")

	CreateVerificationResp, err := client.VerifyV2.CreateVerification(*CreateServiceResp.Sid, _CreateVerificationParams)
	if err != nil {
		fmt.Println("CreateVerificationResp: ", CreateVerificationResp)
		return "Error", "Error"
	}

	return *CreateVerificationResp.Sid, *CreateServiceResp.Sid
}

func VerifyOtp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env.local file")
	}

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	var requestBody struct {
		Phone                string `json:"phone"`
		Otp                  string `json:"otp"`
		VerificationResponse string `json:"verificationResponse"`
		ServiceResponseParam string `json:"serviceResponseParam"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		response := map[string]interface{}{
			"message": "Unable to verify otp",
		}
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		json.NewEncoder(w).Encode(response)

		log.Fatal("Error decoding request body: ", err)
	}

	phone := requestBody.Phone
	otp := requestBody.Otp
	verificationResponse := requestBody.VerificationResponse
	serviceResponseParam := requestBody.ServiceResponseParam
	if phone != "" && otp != "" {
		resp := isCorrectOtp(phone, otp, accountSid, authToken, verificationResponse, serviceResponseParam)
		if resp == "approved" {
			jwtToken := generateJWT(phone)
			response := map[string]interface{}{
				"jwtToken": jwtToken,
				"message":  "Otp verified successfully",
			}
			w.WriteHeader(http.StatusOK) // Set the HTTP status code
			json.NewEncoder(w).Encode(response)
		} else {
			response := map[string]interface{}{
				"message": "Unable to verify otp",
			}
			w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
			json.NewEncoder(w).Encode(response)
		}
	} else {
		response := map[string]interface{}{
			"message": "Unable to verify otp",
		}
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		json.NewEncoder(w).Encode(response)
	}

}

func isCorrectOtp(phone string, otp string, accountSid string, authToken string, verificationResponse string, serviceResponseParam string) string {

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   accountSid,
		Password:   authToken,
		AccountSid: accountSid,
	})

	_CreateVerificationCheckParams := &verify.CreateVerificationCheckParams{}
	_CreateVerificationCheckParams.SetTo(phone)
	_CreateVerificationCheckParams.SetCode(otp)
	_CreateVerificationCheckParams.SetVerificationSid(verificationResponse)

	CreateVerificationCheckResp, err := client.VerifyV2.CreateVerificationCheck(serviceResponseParam, _CreateVerificationCheckParams)
	if err != nil {
		return "Error"
	}
	return *CreateVerificationCheckResp.Status

}

func generateJWT(phone string) string {

	godotenv.Load(".env.local")

	// Define the secret key used to sign the token

	secretKey := os.Getenv("JWT_SECRET")

	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = phone
	claims["exp"] = time.Now().Add(time.Hour * 4).Unix() // Token expiration time (1 hour from now)

	// Generate the JWT string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println("Error generating token:", err)
		return "Error"
	}

	// Print the JWT string
	return tokenString
}

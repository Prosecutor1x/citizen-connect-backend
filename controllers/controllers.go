package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Prosecutor1x/citizen-connect-frontend/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://mukeshkuiry:12345mukesh@cluster0.duhyfwm.mongodb.net/"
const dbName = "problems-list"
const colName = "problem"

var collection *mongo.Collection

func init() {
	clientOption := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB successfully")

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection instance is ready")
}

func insertOneProblem(problem model.ProblemData) {
	inserted, err := collection.InsertOne(context.Background(), problem)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", inserted.InsertedID)
}

func AddDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var problem model.ProblemData
	_ = json.NewDecoder(r.Body).Decode(&problem)

	insertOneProblem(problem)
	json.NewEncoder(w).Encode(problem)

}

func deleteOneIssue(issueId string) {
	id, _ := primitive.ObjectIDFromHex(issueId)
	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted Count: ", result.DeletedCount)
}

func DeleteDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	issueId := params["id"]
	deleteOneIssue(issueId)
	json.NewEncoder(w).Encode(issueId)

}

func getAllIssue() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	var issues []primitive.M
	for cur.Next(context.Background()) {
		var issue bson.M
		err := cur.Decode(&issue)
		if err != nil {
			log.Fatal(err)
		}
		issues = append(issues, issue)
	}
	return issues
}

func FetchAllDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	json.NewEncoder(w).Encode(getAllIssue())
}

func FetchSingleDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	params := mux.Vars(r)
	issueId := params["id"]
	id, _ := primitive.ObjectIDFromHex(issueId)
	filter := bson.M{"_id": id}
	var issue bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&issue)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(issue)
}

func updateOneIssue(issueId string, problem model.ProblemData) {
	id, _ := primitive.ObjectIDFromHex(issueId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": problem}
	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified Count: ", result.ModifiedCount)
}

func UpdateDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	var problem model.ProblemData
	_ = json.NewDecoder(r.Body).Decode(&problem)
	params := mux.Vars(r)
	issueId := params["id"]
	fmt.Println("Issue ID: ", issueId)
	fmt.Println("Problem: ", problem)
	updateOneIssue(issueId, problem)
	json.NewEncoder(w).Encode(issueId)
}

func SendOtp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	err := godotenv.Load(".env")
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
		log.Fatal("Error loading .env file")
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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
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

	godotenv.Load(".env")

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

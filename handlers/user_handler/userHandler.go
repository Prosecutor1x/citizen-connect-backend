package user_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Prosecutor1x/citizen-connect-frontend/database"
	"github.com/Prosecutor1x/citizen-connect-frontend/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const dbName = "user-list"
const colName = "users"

var collection *mongo.Collection

func init() {

	client, err := database.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection instance is ready")
}

func createUser(user model.UserData) {
	inserted, err := collection.InsertOne(context.Background(), user)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", inserted.InsertedID)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var user model.UserData
	_ = json.NewDecoder(r.Body).Decode(&user)

	createUser(user)
	json.NewEncoder(w).Encode(user)

}

func CheckUserExist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var phone model.Phone
	_ = json.NewDecoder(r.Body).Decode(&phone)
	filter := bson.M{"userphone": phone.Phone}

	var result bson.M

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	fmt.Println("result", result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No documents found")
			// send 500 error
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("No documents found")
		} else {
			json.NewEncoder(w).Encode("Error decoding")
		}
	} else {
		json.NewEncoder(w).Encode(result)
	}

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	params := mux.Vars(r)
	userId := params["id"]
	if userId == "" {
		response := map[string]interface{}{
			"message": "Please provide correct user id",
		}
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		json.NewEncoder(w).Encode(response)
		return

	}
	id, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": id}
	var userData bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&userData)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(userData)

}

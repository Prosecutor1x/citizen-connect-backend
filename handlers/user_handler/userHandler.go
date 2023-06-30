package user_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Prosecutor1x/citizen-connect-frontend/database"
	"github.com/Prosecutor1x/citizen-connect-frontend/model"
	"go.mongodb.org/mongo-driver/bson"
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

	var result model.UserData

	err := collection.FindOne(context.Background(), filter).Decode(&result)

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

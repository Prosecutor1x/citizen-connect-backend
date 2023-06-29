package user_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Prosecutor1x/citizen-connect-frontend/database"
	"github.com/Prosecutor1x/citizen-connect-frontend/model"
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

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Prosecutor1x/citizen-connect-frontend/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
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

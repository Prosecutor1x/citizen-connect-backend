package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Prosecutor1x/citizen-connect-frontend/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

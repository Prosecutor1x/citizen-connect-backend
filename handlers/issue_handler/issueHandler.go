package issue_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/Prosecutor1x/citizen-connect-frontend/database"
	"github.com/Prosecutor1x/citizen-connect-frontend/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const dbName = "problems_list"
const colName = "problem"
var nan = math.NaN()


var collection *mongo.Collection

func init() {

	client, err := database.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection instance is ready")
}

func AddOneIssue(problem model.ProblemData) {
	inserted, err := collection.InsertOne(context.Background(), problem)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", inserted.InsertedID)
}

func CreateIssueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var problem model.ProblemData
	_ = json.NewDecoder(r.Body).Decode(&problem)

	if problem.IssueTitle == "" ||
		problem.IssueDescription == "" ||
		problem.IssueLocation.Lat == nan ||
		problem.IssueLocation.Long == nan ||
		problem.IssueProgress == "" ||
		problem.IssueRaiserId == "" ||
		problem.IssueDate == "" {
		response := map[string]interface{}{
			"message": "Please fill all the fields",
		}
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		json.NewEncoder(w).Encode(response)
		return
	}

	AddOneIssue(problem)
	json.NewEncoder(w).Encode(problem)

}

func DeleteOneIssue(issueId string) {
	id, _ := primitive.ObjectIDFromHex(issueId)
	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted Count: ", result.DeletedCount)
}

func DeleteIssueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	issueId := params["id"]
	if issueId == "" {
		response := map[string]interface{}{
			"message": "Please provide correct issue id",
		}
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		json.NewEncoder(w).Encode(response)
		return

	}
	DeleteOneIssue(issueId)
	json.NewEncoder(w).Encode(issueId)

}

func GetAllIssue() []primitive.M {
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

func FetchAllIssueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	json.NewEncoder(w).Encode(GetAllIssue())
}

func FetchSingleIssueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	params := mux.Vars(r)
	issueId := params["id"]
	if issueId == "" {
		response := map[string]interface{}{
			"message": "Please provide correct issue id",
		}
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		json.NewEncoder(w).Encode(response)
		return

	}
	id, _ := primitive.ObjectIDFromHex(issueId)
	filter := bson.M{"_id": id}
	var issue bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&issue)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(issue)
}

func UpdateOneIssue(issueId string, problem model.ProblemData) {
	id, _ := primitive.ObjectIDFromHex(issueId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": problem}
	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified Count: ", result.ModifiedCount)
}

func UpdateIssueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	var problem model.ProblemData
	_ = json.NewDecoder(r.Body).Decode(&problem)
	params := mux.Vars(r)
	issueId := params["id"]
	if issueId == "" {
		response := map[string]interface{}{
			"message": "Please provide correct issue id",
		}
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		json.NewEncoder(w).Encode(response)
		return
	}
	if problem.IssueTitle == "" ||
		problem.IssueDescription == "" ||
		problem.IssueLocation.Long == nan ||
		problem.IssueLocation.Lat == nan ||
		problem.IssueProgress == "" ||
		problem.IssueRaiserId == "" ||
		problem.IssueDate == "" {
		response := map[string]interface{}{
			"message": "Please fill all the fields",
		}
		w.WriteHeader(http.StatusInternalServerError) // Set the HTTP status code
		json.NewEncoder(w).Encode(response)
		return
	}
	UpdateOneIssue(issueId, problem)
	json.NewEncoder(w).Encode(issueId)
}

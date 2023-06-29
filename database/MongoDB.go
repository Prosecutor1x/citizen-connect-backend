package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Client, error) {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal(err)
	}
	connectionString := os.Getenv("MONGO_URL")
	clientOption := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB successfully")
	return client, nil
}

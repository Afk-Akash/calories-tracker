package main

import (
    "log"
	"context"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "calorie-tracker/handlers" // Adjust this import path as necessary
)

func main() {
    // Create a new Fiber app
    app := fiber.New()

    // MongoDB connection setup
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Replace with your MongoDB URI
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatalf("Could not connect to MongoDB: %v", err)
    }

    // Ensure the connection is established
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatalf("Could not ping MongoDB: %v", err)
    }
    log.Println("Connected to MongoDB!")

    // Set up the user collection
    handlers.SetUpUserCollection(client)

    // Define routes
    app.Post("/register", handlers.Register)
    app.Post("/login", handlers.Login)

    // Start the server
    log.Fatal(app.Listen(":3000")) 
}

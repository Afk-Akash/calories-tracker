package main

import (
	"context"
	"log"

	"calorie-tracker/handlers" // Adjust this import path as necessary
	"calorie-tracker/middlewares"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	handlers.SetUpIngredientCollection(client)
    handlers.SetUpMealCollection(client)

    // Define routes
    app.Post("/register", handlers.Register)
    app.Post("/login", handlers.Login)

	api := app.Group("/api", middlewares.AuthMiddleware)
    api.Post("/ingredient", handlers.CreateIngredient)
	api.Put("/update/ingredient/:id", handlers.UpdateIngredient)
	api.Delete("/delete/ingredient/:id", handlers.DeleteIngredient)
    api.Post("/addMeal", handlers.CreateMeal)


    // Start the server
    log.Fatal(app.Listen(":3000")) 
}

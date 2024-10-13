package main

import (
	"context"
	"fmt"
	"log"

	"calorie-tracker/handlers" // Adjust this import path as necessary
	"calorie-tracker/middlewares"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
<<<<<<< HEAD
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
=======
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
    log.Println("creating userDatabase and collection")
    handlers.SetUpUserCollection(client)
    log.Println("creating database Ingredient and collection")
	handlers.SetUpIngredientCollection(client)
    log.Println("creating database meal and collection")
    handlers.SetUpMealCollection(client)
    log.Println("creating database daily log tracker and collection")
    handlers.SetUpDailyLogCollection(client)
>>>>>>> adcc75b6fe9a439b1a4bf6394cbd44fe9d82530c

	// Define routes
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	api := app.Group("/api", middlewares.AuthMiddleware)
	api.Post("/ingredient", handlers.CreateIngredient)
	api.Put("/update/ingredient/:id", handlers.UpdateIngredient)
	api.Delete("/delete/ingredient/:id", handlers.DeleteIngredient)
<<<<<<< HEAD
	api.Post("/addMeal", handlers.CreateMeal)
	api.Get("/GetUserIngredients", handlers.GetUserIngredients)

	// Start the server
	log.Fatal(app.Listen(":3000"))
=======
    api.Post("/createMeal", handlers.CreateMeal)
    api.Post("/addMeal", handlers.CreateDailyLog)


    // Start the server
    fmt.Println("starting the server....")
    log.Fatal(app.Listen(":3000")) 
>>>>>>> adcc75b6fe9a439b1a4bf6394cbd44fe9d82530c
}

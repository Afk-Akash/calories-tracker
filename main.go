package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"calorie-tracker/handlers" // Adjust this import path as necessary
	"calorie-tracker/middlewares"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Create a new Fiber app
	app := fiber.New()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	log.Println("trying to connect with mongo with url : ", mongoURI)
	// MongoDB connection setup
	clientOptions := options.Client().ApplyURI(mongoURI) // Replace with your MongoDB URI
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

	// Define routes
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	api := app.Group("/api", middlewares.AuthMiddleware)
	api.Post("/ingredient", handlers.CreateIngredient)
	api.Post("/update/ingredient", handlers.UpdateIngredient)
	api.Delete("/delete/ingredient", handlers.DeleteIngredient)
	api.Post("/createMeal", handlers.CreateMeal)
	api.Post("/update/meal", handlers.DeleteMeal)
	api.Post("/addMeal", handlers.CreateDailyLog)
	api.Get("/GetUserIngredients", handlers.GetUserIngredients)
	api.Get("/getDailyLog", handlers.GetDailyLogs)

	// Start the server
	fmt.Println("starting the server....")
	log.Fatal(app.Listen(":3000"))
}

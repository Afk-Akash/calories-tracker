package handlers

import (
	"calorie-tracker/models"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var mealCollection *mongo.Collection

func SetUpMealCollection(client *mongo.Client) {
    mealCollection = client.Database("calorieTracker").Collection("meals", &options.CollectionOptions{
        WriteConcern: writeconcern.New(writeconcern.W(1)),
    })

    // Create a unique index on the "name" field
    indexModel := mongo.IndexModel{
        Keys:    map[string]interface{}{"name": 1}, // Ascending index on "name"
        Options: options.Index().SetUnique(true),  // Ensure the name is unique
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Create the index
    if _, err := mealCollection.Indexes().CreateOne(ctx, indexModel); err != nil {
        log.Fatalf("Failed to create index: %v", err)
    }

    log.Println("Meal collection set up with unique index on 'name' field")
}

func CreateMeal(c *fiber.Ctx) error {
    user := c.Locals("user").(map[string]interface{})
    userID := user["user_id"].(string)

    var meal models.Meal
    if err := c.BodyParser(&meal); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    var totalCalories, totalFat, totalCarbs, totalProtien float64
    for i := 0; i < len(meal.Ingredients); i++ {
        totalCalories = totalCalories + meal.Ingredients[i].Calories
        totalProtien  = totalProtien + meal.Ingredients[i].Protein
        totalCarbs    = totalCarbs + meal.Ingredients[i].Carbs
        totalFat      = totalFat + meal.Ingredients[i].Fat
    }

    meal.TotalCalories = totalCalories
    meal.TotalProtien  = totalProtien
    meal.TotalCarbs    = totalCarbs
    meal.TotalFat      = totalFat

    objectId, _ := primitive.ObjectIDFromHex(userID)
    meal.UserID = objectId
    meal.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

    // Insert the ingredient into MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := mealCollection.InsertOne(ctx, meal)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create meal",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "meal created",
    })
}

package handlers

import (
	"calorie-tracker/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dailyLogCollection *mongo.Collection

func SetUpDailyLogCollection(client *mongo.Client) {
    dailyLogCollection = client.Database("calorieTracker").Collection("daily_logs")
    mealCollection     = client.Database("calorieTracker").Collection("meals")
}

func CreateDailyLog(c *fiber.Ctx) error {
    user := c.Locals("user").(map[string]interface{})
	userID := user["user_id"].(string)

	// Parse the request body into a Meal object
	var meal models.Meal
	if err := c.BodyParser(&meal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
    mealJSON, err := json.MarshalIndent(meal, "", "  ") // Pretty print with indentation
	if err != nil {
		log.Printf("Error converting meal to JSON: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process meal data",
		})
	}

	// Print the JSON to console
	log.Println(string(mealJSON))

	// Convert userID to MongoDB ObjectID
	objectID, _ := primitive.ObjectIDFromHex(userID)

	// Get the current date (formatted as YYYY-MM-DD) to use as a unique identifier for the day
	today := time.Now()
    todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

    

	// Create the query filter to find the document by user_id and today's date
	filter := bson.M{
		"user_id": objectID,
		"date":    primitive.NewDateTimeFromTime(todayDate),
	}

	// Update operation: Add the meal to the "meals" array and increment the total macros


	opts := options.Update().SetUpsert(true)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

    var existingMeal models.Meal
    mealCollection.FindOne(ctx, bson.M{"name": meal.Name, "user_id": objectID}).Decode(&existingMeal)
    vv , _:= json.MarshalIndent(meal, "", "  ")
    log.Println(string(vv))
    update := bson.M{
		"$push": bson.M{"meals": existingMeal}, // Push the new meal to the meals array
		"$inc": bson.M{
			"total_calories": existingMeal.TotalCalories,
			"total_protein":  existingMeal.TotalProtien,
			"total_carbs":    existingMeal.TotalCarbs,
			"total_fat":      existingMeal.TotalFat,
		},
		"$setOnInsert": bson.M{ // Set these values only if the document is newly created
			"user_id":  objectID,
			"date":     todayDate,
			"created_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
    
	_, err = dailyLogCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add meal, seems like database error...Please try again after some time",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Meal added successfully",
	})
}

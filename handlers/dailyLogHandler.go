package handlers

import (
	"calorie-tracker/models"
	"context"
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
	update := bson.M{
		"$push": bson.M{"meals": meal}, // Push the new meal to the meals array
		"$inc": bson.M{
			"total_calories": meal.TotalCalories,
			"total_protein":  meal.TotalProtien,
			"total_carbs":    meal.TotalCarbs,
			"total_fat":      meal.TotalFat,
		},
		"$setOnInsert": bson.M{ // Set these values only if the document is newly created
			"user_id":  objectID,
			"date":     primitive.NewDateTimeFromTime(time.Now()),
			"created_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	// Use the upsert option to create the document if it doesn't exist
	opts := options.Update().SetUpsert(true)

	// Perform the update operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := dailyLogCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add meal, seems like database error...Please try again after some time",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Meal added successfully",
	})
}

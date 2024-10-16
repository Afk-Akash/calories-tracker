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

	objectID, _ := primitive.ObjectIDFromHex(userID)

	today := time.Now()
    todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

    

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
 
    update := bson.M{
		"$push": bson.M{"meals": existingMeal}, 
		"$inc": bson.M{
			"total_calories": existingMeal.TotalCalories,
			"total_protien":  existingMeal.TotalProtien,
			"total_carbs":    existingMeal.TotalCarbs,
			"total_fat":      existingMeal.TotalFat,
		},
		"$setOnInsert": bson.M{
			"user_id":  objectID,
			"date":     todayDate,
			"created_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
    
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

func GetDailyLogs(c *fiber.Ctx) error {
	user := c.Locals("user").(map[string]interface{})
	userID := user["user_id"].(string)

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get the current date (midnight) to filter logs by date
	today := time.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	// Create the query filter for user_id and date
	filter := bson.M{
		"user_id": objectID,
		"date":    primitive.NewDateTimeFromTime(todayDate),
	}

	// Find the daily log in the database
	var dailyLog models.DailyLog
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = dailyLogCollection.FindOne(ctx, filter).Decode(&dailyLog)
	vv , _:= json.MarshalIndent(dailyLog, "", "  ")
    log.Println(string(vv))
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No log found for today",
		})
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve daily log, seems like db error",
		})
	}

	// Return the daily log as JSON
	return c.Status(fiber.StatusOK).JSON(dailyLog)
}

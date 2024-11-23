package handlers

import (
	"calorie-tracker/models"
	Request "calorie-tracker/request"
	"context"
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
	mealCollection = client.Database("calorieTracker").Collection("meals")
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
	err := mealCollection.FindOne(ctx, bson.M{"name": meal.Name, "user_id": objectID}).Decode(&existingMeal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No meal found with the specified name",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An error occurred while fetching the meal",
		})
	}

	existingMeal.LogID = primitive.NewObjectID()

	update := bson.M{
		"$push": bson.M{"meals": existingMeal},
		"$inc": bson.M{
			"total_calories": existingMeal.TotalCalories,
			"total_protein":  existingMeal.TotalProtein,
			"total_carbs":    existingMeal.TotalCarbs,
			"total_fat":      existingMeal.TotalFat,
		},
		"$setOnInsert": bson.M{
			"user_id":    objectID,
			"date":       todayDate,
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

func DeleteMealFromDailyLog(c *fiber.Ctx) error {
	user := c.Locals("user").(map[string]interface{})
	userID := user["user_id"].(string)

	var req Request.DeleteMealFromDailyLog

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
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

	// Prepare the update to pull the meal and adjust macros
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the meal to delete in the `meals` array
	var dailyLog models.DailyLog
	err = dailyLogCollection.FindOne(ctx, filter).Decode(&dailyLog)
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No daily log found for today",
		})
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve daily log",
		})
	}

	// Find the meal to delete
	var mealToDelete *models.Meal
	for _, meal := range dailyLog.Meals {
		if meal.LogID == req.LogID {
			mealToDelete = &meal
			break
		}
	}

	if mealToDelete == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Meal not found in daily log",
		})
	}

	// Prepare the update to remove the meal and adjust the macros
	update := bson.M{
		"$pull": bson.M{
			"meals": bson.M{"_id": req.LogID},
		},
		"$inc": bson.M{
			"total_calories": -mealToDelete.TotalCalories,
			"total_protein":  -mealToDelete.TotalProtein,
			"total_carbs":    -mealToDelete.TotalCarbs,
			"total_fat":      -mealToDelete.TotalFat,
		},
	}

	// Update the daily log
	_, err = dailyLogCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Database error during update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update daily log",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Meal deleted successfully and macros adjusted",
	})
}

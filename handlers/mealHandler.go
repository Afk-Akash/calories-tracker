package handlers

import (
	"calorie-tracker/models"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
		Options: options.Index().SetUnique(true),   // Ensure the name is unique
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

	if meal.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Field 'name' is required",
		})
	}

	if meal.Ingredients == nil || len(meal.Ingredients) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Field 'ingredients' is required and cannot be empty",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var totalCalories, totalFat, totalCarbs, totalProtein float64
	for i := 0; i < len(meal.Ingredients); i++ {
		var dbIngredient models.Ingredient
		filter := bson.M{"name": meal.Ingredients[i].Name}
		err := ingredientCollection.FindOne(ctx, filter).Decode(&dbIngredient)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Ingredient not found in database",
					"missing": meal.Ingredients[i].Name,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error while validating ingredient",
			})
		}

		totalCalories = totalCalories + dbIngredient.Calories
		meal.Ingredients[i].Calories = dbIngredient.Calories

		totalProtein = totalProtein + dbIngredient.Protein
		meal.Ingredients[i].Protein = dbIngredient.Protein

		totalCarbs = totalCarbs + dbIngredient.Carbs
		meal.Ingredients[i].Carbs = dbIngredient.Carbs

		totalFat = totalFat + dbIngredient.Fat
		meal.Ingredients[i].Fat = dbIngredient.Fat

	}

	meal.TotalCalories = totalCalories
	meal.TotalProtein = totalProtein
	meal.TotalCarbs = totalCarbs
	meal.TotalFat = totalFat

	objectId, _ := primitive.ObjectIDFromHex(userID)
	meal.UserID = objectId
	meal.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Insert the ingredient into MongoDB

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

func DeleteMeal(c *fiber.Ctx) error {
	user := c.Locals("user").(map[string]interface{})
	userID := user["user_id"].(string)

	var meal models.Meal
	if err := c.BodyParser(&meal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if meal.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Meal name is required",
		})
	}

	// Convert userID to MongoDB ObjectID
	objectId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Create a MongoDB filter to delete the meal based on user ID and meal name
	filter := bson.M{
		"user_id": objectId,
		"name":    meal.Name,
	}

	// Set a timeout for the MongoDB operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete the meal from MongoDB
	result, err := mealCollection.DeleteOne(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete meal",
		})
	}

	// Check if a meal was actually deleted
	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Meal not found",
		})
	}

	// Respond with success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Meal deleted successfully",
	})
}

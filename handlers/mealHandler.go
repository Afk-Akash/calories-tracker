package handlers

import (
    "calorie-tracker/models"
    "context"
    "time"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var mealCollection *mongo.Collection

func SetUpMealCollection(client *mongo.Client) {
    mealCollection = client.Database("calorieTracker").Collection("meals")
}

func CreateMeal(c *fiber.Ctx) error {
    var meal models.Meal
    if err := c.BodyParser(&meal); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
    }

    meal.ID = primitive.NewObjectID()
    meal.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

    // Calculate total macros and calories
    totalMacros := models.Macros{}
    for _, ingredient := range meal.Ingredients {
        totalMacros.Protein += ingredient.Macros.Protein
        totalMacros.Carbs += ingredient.Macros.Carbs
        totalMacros.Fats += ingredient.Macros.Fats
        meal.TotalCalories += ingredient.Calories
    }
    meal.TotalMacros = totalMacros

    _, err := mealCollection.InsertOne(context.Background(), meal)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to insert meal"})
    }

    return c.Status(fiber.StatusOK).JSON(meal)
}

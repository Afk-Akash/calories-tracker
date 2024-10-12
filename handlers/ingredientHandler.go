package handlers

import (
    "calorie-tracker/models"
    "context"
    "time"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var ingredientCollection *mongo.Collection

func SetUpIngredientCollection(client *mongo.Client) {
    ingredientCollection = client.Database("calorieTracker").Collection("ingredients")
}

func CreateIngredient(c *fiber.Ctx) error {
    var ingredient models.Ingredient
    if err := c.BodyParser(&ingredient); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
    }

    ingredient.ID = primitive.NewObjectID()
    ingredient.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

    _, err := ingredientCollection.InsertOne(context.Background(), ingredient)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to insert ingredient"})
    }

    return c.Status(fiber.StatusOK).JSON(ingredient)
}

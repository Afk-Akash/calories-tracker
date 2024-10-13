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

var ingredientCollection *mongo.Collection

func SetUpIngredientCollection(client *mongo.Client) {
    ingredientCollection = client.Database("calorieTracker").Collection("ingredients", &options.CollectionOptions{
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

func CreateIngredient(c *fiber.Ctx) error {
    user := c.Locals("user").(map[string]interface{})
    userID := user["user_id"].(string)

    // Parse request body into the Ingredient struct
    var ingredient models.Ingredient
    if err := c.BodyParser(&ingredient); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    // Add the user ID to the ingredient
    objectId, _ := primitive.ObjectIDFromHex(userID)
    ingredient.UserID = objectId
    ingredient.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

    // Insert the ingredient into MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := ingredientCollection.InsertOne(ctx, ingredient)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to add ingredient",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Ingredient added successfully",
    })
}


func UpdateIngredient(c *fiber.Ctx) error {
    // Get the ingredient ID from the URL parameters
    ingredientID := c.Params("id")
    if ingredientID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Missing ingredient ID",
        })
    }

    user := c.Locals("user").(map[string]interface{})
    userID := user["user_id"].(string)

    objectID, err := primitive.ObjectIDFromHex(ingredientID)
    userId, _ := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid ingredient ID",
        })
    }

    var existingIngredient, ingredient models.Ingredient

    err = ingredientCollection.FindOne(context.Background(), bson.M{"_id": objectID, "user_id": userId}).Decode(&existingIngredient)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Ingredient not found or does not belong to the user",
        })
    }

    if err := c.BodyParser(&ingredient); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Create an update document
    update := bson.M{
        "$set": bson.M{
            "protein":  ingredient.Protein,
            "carbs":    ingredient.Carbs,
            "fat":     ingredient.Fat,
            "calories": ingredient.Calories,
        },
    }

    // Update the ingredient in MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := ingredientCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update ingredient",
        })
    }

    if result.ModifiedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "couldn't modify the ingredient, seems like it is already updated",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Ingredient updated successfully",
    })
}


func DeleteIngredient(c *fiber.Ctx) error {
    // Get the user from the context (optional)
    user := c.Locals("user").(map[string]interface{})
    userID := user["user_id"].(string)

    // Get the ingredient ID from the URL parameters
    ingredientID := c.Params("id")
    
    // Create a context with timeout for the database operation
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Convert the ingredient ID to ObjectID
    objectId, err := primitive.ObjectIDFromHex(ingredientID)
    userId, _ := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid ingredient ID",
        })
    }

    result := ingredientCollection.FindOne(ctx, bson.M{"_id": objectId, "user_id": userId})
    if result.Err() != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Ingredient not found",
        })
    }

    _, err = ingredientCollection.DeleteOne(ctx, bson.M{"_id": objectId, "user_id": userId})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to delete ingredient",
        })
    }

    return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
        "message": "Ingredient deleted successfully",
    })
}

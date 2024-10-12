package handlers

import (
    "calorie-tracker/models"
    "context"
    "time"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var dailyLogCollection *mongo.Collection

func SetUpDailyLogCollection(client *mongo.Client) {
    dailyLogCollection = client.Database("calorieTracker").Collection("daily_logs")
}

func CreateDailyLog(c *fiber.Ctx) error {
    var dailyLog models.DailyLog
    if err := c.BodyParser(&dailyLog); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
    }

    dailyLog.ID = primitive.NewObjectID()
    dailyLog.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

    _, err := dailyLogCollection.InsertOne(context.Background(), dailyLog)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to insert daily log"})
    }

    return c.Status(fiber.StatusOK).JSON(dailyLog)
}

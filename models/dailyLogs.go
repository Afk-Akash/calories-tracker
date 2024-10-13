package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DailyLog struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
    Date         primitive.DateTime `bson:"date" json:"date"`
    Meals        []Meal             `bson:"meals" json:"meals"`
    TotalCalories float64           `bson:"total_calories" json:"total_calories"`
    CreatedAt    primitive.DateTime `bson:"created_at" json:"created_at"`
}

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Meal struct {
    ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
    Name          string             `bson:"name" json:"name"`
    Ingredients   []Ingredient       `bson:"ingredients" json:"ingredients"`
    TotalCalories float64            `bson:"total_calories" json:"total_calories"`
    TotalProtien  float64            `bson:"total_protien" json:"total_protien"`
    TotalCarbs    float64            `bson:"total_carbs" json:"total_carbs"`
    TotalFat      float64            `bson:"total_fat" json:"total_fat"`
    CreatedAt     primitive.DateTime `bson:"created_at" json:"created_at"`
}

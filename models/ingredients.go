package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Macros struct {
    Protein float64 `bson:"protein" json:"protein"`
    Carbs   float64 `bson:"carbs" json:"carbs"`
    Fats    float64 `bson:"fats" json:"fats"`
}

type Ingredient struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
    Name      string             `bson:"name" json:"name"`
    Macros    Macros             `bson:"macros" json:"macros"`
    Calories  float64            `bson:"calories" json:"calories"`
    CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}

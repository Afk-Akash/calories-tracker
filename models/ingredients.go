package models


import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Ingredient struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
    Name      string             `bson:"name" json:"name"`
    Protein float64              `bson:"protein" json:"protein"`
    Carbs   float64              `bson:"carbs" json:"carbs"`
    Fat    float64               `bson:"fat" json:"fat"`
    Calories  float64            `bson:"calories" json:"calories"`
    CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}

func NewIngredient(userID primitive.ObjectID, name string, protein, carbs, fat, calories float64) *Ingredient {
    return &Ingredient{
        ID:        primitive.NewObjectID(), // Assign a new ObjectID
        UserID:    userID,
        Name:      name,
        Protein:   protein,
        Carbs:     carbs,
        Fat:       fat,
        Calories:  calories,
        CreatedAt: primitive.NewDateTimeFromTime(time.Now()), // Set created at to now
    }
}
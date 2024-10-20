package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Username string             `json:"username" bson:"username"`
    Email    string             `json:"email" bson:"email"`
    Password string             `json:"-" bson:"password"` // Omit from JSON response for security
}

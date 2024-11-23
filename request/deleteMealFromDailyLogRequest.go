package Request

import "go.mongodb.org/mongo-driver/bson/primitive"

type DeleteMealFromDailyLog struct {
	LogID      primitive.ObjectID `bson:"logID,omitempty" json:"logID,omitempty"`
}
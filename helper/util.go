package helper

import (
	"context"
	"emobackend/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type SystemPrompt struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Text string             `bson:"text" json:"text"`
}


func GetPromptByID(id string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	var prompt SystemPrompt
	err = config.DB.Collection("ai_prompts").
		FindOne(context.TODO(), bson.M{"_id": objectID}).
		Decode(&prompt)
	if err != nil {
		return "", err
	}

	return prompt.Text, nil
}

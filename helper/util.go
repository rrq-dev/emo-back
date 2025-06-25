package helper

import (
	"context"
	"emobackend/config"

	"gopkg.in/mgo.v2/bson"
)

func GetSystemPrompt() (string, error) {
	collection := config.DB.Collection("ai_prompts")

	var prompt struct {
		Text string `bson:"text"`
	}

	err := collection.FindOne(context.TODO(), bson.M{"_id": "default_reflection"}).Decode(&prompt)
	if err != nil {
		return "", err
	}

	return prompt.Text, nil
}

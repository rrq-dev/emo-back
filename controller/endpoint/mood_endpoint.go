package controller

import (
	"context"
	"emobackend/config"
	gene "emobackend/helper"
	"emobackend/model"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func GetAllMoodsData(c *fiber.Ctx) error {
	collection := config.DB.Collection("heatmaps")
	ctx := context.Background()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var moods []model.MoodReflection
	if err := cursor.All(ctx, &moods); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	return c.JSON(moods)
}

func SubmitMoods(c *fiber.Ctx) error {
	var input model.MoodInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input tidak valid",
			"error":   err.Error(),
		})
	}

	var userID, userName string

	if input.IsAnonymous {
		userID = "anon-" + gene.GenerateUserID()
		userName = "Anonymous"
	} else {
		userIDHex := c.Locals("user_id").(string)
		objID, err := primitive.ObjectIDFromHex(userIDHex)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID from token",
			})
		}

		var userData model.User
		err = config.DB.Collection("users").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&userData)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User tidak ditemukan",
				"error":   err.Error(),
			})
		}

		userID = "user-" + userData.ID.Hex()
		userName = userData.Name
	}

	moods := model.MoodReflection{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		UserName:   userName,
		Mood:       input.Mood,
		Reflection: input.Message,
		Timestamp:  time.Now(),
	}

	go func() {
		if input.Message != "" {
			err := callGeminiAndSaveReflection(userID, input.Message, input.IsAnonymous)
			if err != nil {
				fmt.Println("Gagal menyimpan refleksi Gemini:", err)
			}
		}
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil menyimpan refleksi mood",
		"data":    moods,
	})
}

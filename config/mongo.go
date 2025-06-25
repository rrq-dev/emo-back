package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Database
var DBName = "emodata"

func MongoConnect() {
	mongoString := os.Getenv("MONGODB_URL")
	if mongoString == "" {
		log.Fatal("❌ MongoConnect: MONGODB_URL tidak ditemukan di environment variable")
	}

	// Setup client options
	clientOpts := options.Client().ApplyURI(mongoString).
		SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatalf("❌ MongoConnect: Gagal koneksi ke MongoDB: %v", err)
	}

	// Ping database untuk pastikan bisa connect
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatalf("❌ MongoConnect: Gagal ping MongoDB: %v", err)
	}

	DB = client.Database(DBName)
	fmt.Println("✅ MongoConnect: Berhasil terhubung ke MongoDB 🎉")
}

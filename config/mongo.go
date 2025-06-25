package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Database
var MongoString = os.Getenv("MONGO_URL")
var DBName = "emodata"

func MongoConnect() {
	mongoString := os.Getenv("MONGO_URL")
	if mongoString == "" {
		fmt.Println("MongoConnect: MONGO_URL tidak ditemukan di .env")
		return
	}

	clientOpts := options.Client().ApplyURI(mongoString)

	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		fmt.Println("MongoConnect: gagal koneksi:", err)
		return
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println("MongoConnect: gagal ping MongoDB:", err)
		return
	}

	fmt.Println("MongoConnect: berhasil terhubung ke MongoDB")
	DB = client.Database(DBName)
}

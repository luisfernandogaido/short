package model

import (
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	db *mongo.Database
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Ini(mongoURI string) error {
	client, err := mongo.Connect(nil, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("model ini: %w", err)
	}
	err = client.Ping(nil, readpref.Primary())
	if err != nil {
		return fmt.Errorf("model ini: %w", err)
	}
	db = client.Database("short")
	return nil
}

package model

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	db            *mongo.Database
	rdb           *redis.Client
	ctx           = context.Background()
	ErrNotFound   = errors.New("not found")
	ErrDuplicated = errors.New("duplicated")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Ini(mongoURI string, redisURL string) error {
	client, err := mongo.Connect(nil, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("model ini: %w", err)
	}
	err = client.Ping(nil, readpref.Primary())
	if err != nil {
		return fmt.Errorf("model ini: %w", err)
	}
	db = client.Database("short")
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("model ini: %w", err)
	}
	rdb = redis.NewClient(opt)
	if err := loadTokens(); err != nil {
		log.Fatal(err)
	}
	return nil
}

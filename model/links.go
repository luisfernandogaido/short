package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	lenHash = 7
)

type Link struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Hash        string             `json:"hash" bson:"hash"`
	Destination string             `json:"destination" bson:"destination"`
	User        string             `json:"user"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

func LinkCreate(dest string, u User) (Link, error) {
	link := Link{
		Destination: dest,
		Hash:        generateHash(),
		User:        u.Name,
		CreatedAt:   time.Now(),
	}
	ior, err := db.Collection("links").InsertOne(nil, link)
	if err != nil {
		return Link{}, fmt.Errorf("link create: %w", err)
	}
	link.Id = ior.InsertedID.(primitive.ObjectID)
	return link, nil
}

func LinkGet(hash string) (Link, error) {
	var link Link
	if err := db.Collection("links").FindOne(nil, bson.D{{"hash", hash}}).Decode(&link); err != nil {
		if err == mongo.ErrNoDocuments {
			return Link{}, fmt.Errorf("model link get: %w", ErrNotFound)
		}
		return Link{}, fmt.Errorf("model link get: %w", err)
	}
	return link, nil
}

func generateHash() string {
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	var hash strings.Builder
	for i := 0; i < lenHash; i++ {
		i := rand.Intn(64)
		hash.WriteString(chars[i : i+1])
	}
	return hash.String()
}

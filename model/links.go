package model

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	lenHash        = 7
	purgePeriod    = 3600
	defaultTtlDays = 3650
)

type Link struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Hash        string             `json:"hash" bson:"hash"`
	Destination string             `json:"destination" bson:"destination"`
	User        string             `json:"user"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt   time.Time          `json:"expires_at" bson:"expires_at"`
}

func LinkCreate(dest string, hash string, ttlDays int, u User) (Link, error) {
	if hash == "" {
		hash = generateHash()
	}
	if ttlDays == 0 {
		ttlDays = defaultTtlDays
	}
	now := time.Now()
	expiresAt := now.Add(time.Hour * 24 * time.Duration(ttlDays))
	link := Link{
		Destination: dest,
		Hash:        hash,
		User:        u.Name,
		CreatedAt:   now,
		ExpiresAt:   expiresAt,
	}
	linkSaved, err := LinkGet(hash)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return Link{}, fmt.Errorf("link create: %w", err)
	}
	if linkSaved.Id == primitive.NilObjectID {
		ior, err := db.Collection("links").InsertOne(nil, link)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return Link{}, fmt.Errorf("link create: %w", ErrDuplicated)
			}
			return Link{}, fmt.Errorf("link create: %w", err)
		}
		link.Id = ior.InsertedID.(primitive.ObjectID)
	} else {
		if linkSaved.User != u.Name {
			return Link{}, fmt.Errorf("link create: %w", ErrDuplicated)
		}
		filter := bson.D{
			{"_id", linkSaved.Id},
		}
		update := bson.D{
			{"$set", bson.D{
				{"destination", dest},
				{"expires_at", expiresAt},
			}},
		}
		if _, err := db.Collection("links").UpdateOne(nil, filter, update); err != nil {
			return Link{}, fmt.Errorf("link create: %w", err)
		}
	}
	go func() {
		db.Collection("links_log").InsertOne(nil, link)
	}()
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

func LinkPurge() {
	for {
		filter := bson.D{
			{"expires_at", bson.D{
				{"$lt", time.Now()},
			}},
		}
		if _, err := db.Collection("links").DeleteMany(nil, filter); err != nil {
			log.Printf("link purge: %v\n", err)
		}
		time.Sleep(time.Second * purgePeriod)
	}
}

func generateHash() string {
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	for {
		var hash strings.Builder
		for i := 0; i < lenHash; i++ {
			i := rand.Intn(64)
			hash.WriteString(chars[i : i+1])
		}
		s := hash.String()
		if !strings.HasSuffix(s, "-") {
			return s
		}
	}
}

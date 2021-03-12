package model

import (
	"fmt"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Token string             `json:"token" bson:"token"`
}

func Users() ([]User, error) {
	cur, err := db.Collection("users").Find(nil, bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("users: %w", err)
	}
	var users []User
	if err := cur.All(nil, &users); err != nil {
		return nil, fmt.Errorf("users: %w", err)
	}
	return users, nil
}

func NewUser(id string) (User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return User{}, fmt.Errorf("new user: %w", err)
	}
	sr := db.Collection("users").FindOne(nil, bson.D{{"_id", oid}})
	var u User
	if err := sr.Decode(&u); err != nil {
		return User{}, fmt.Errorf("new user: %w", err)
	}
	return u, nil
}

func (u *User) Save() error {
	if u.Id == primitive.NilObjectID {
		token, err := generateToken()
		if err != nil {
			return fmt.Errorf("users save: %w", err)
		}
		u.Token = token
		if ior, err := db.Collection("users").InsertOne(nil, u); err != nil {
			return fmt.Errorf("users save: %w", err)
		} else {
			u.Id = ior.InsertedID.(primitive.ObjectID)
		}
		if err := loadTokens(); err != nil {
			return err
		}
		return nil
	}
	if err := db.Collection("users").FindOneAndReplace(nil, bson.D{{"_id", u.Id}}, u).Err(); err != nil {
		return fmt.Errorf("users save: %w", err)
	}
	if err := loadTokens(); err != nil {
		return err
	}
	return nil
}

func (u *User) RegenerateToken() error {
	token, err := generateToken()
	if err != nil {
		return fmt.Errorf("users regeneratetoken: %w", err)
	}
	u.Token = token
	if err := u.Save(); err != nil {
		return fmt.Errorf("users regeneratetoken: %w", err)
	}
	if err := loadTokens(); err != nil {
		return err
	}
	if err := loadTokens(); err != nil {
		return err
	}
	return nil
}

func (u *User) Delete() error {
	if _, err := db.Collection("users").DeleteOne(nil, bson.D{{"_id", u.Id}}); err != nil {
		return fmt.Errorf("users delete: %w", err)
	}
	return nil
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("user save: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}

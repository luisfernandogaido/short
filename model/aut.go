package model

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func Aut(token string) error {
	is, err := rdb.HExists(ctx, "short-users", token).Result()
	if err != nil {
		return fmt.Errorf("model aut: %w", err)
	}
	if !is {
		return fmt.Errorf("model aut: is not member")
	}
	return nil
}

func UserToken(token string) (User, error) {
	b, err := rdb.HGet(ctx, "short-users", token).Bytes()
	if err != nil {
		return User{}, fmt.Errorf("model usertoken: %w", err)
	}
	var u User
	if err := json.Unmarshal(b, &u); err != nil {
		return User{}, fmt.Errorf("model usertoken: %w", err)
	}
	return u, nil
}

func loadTokens() error {
	cur, err := db.Collection("users").Find(nil, bson.D{{}})
	if err != nil {
		return fmt.Errorf("loadtokens: %w", err)
	}
	var users []User
	if err := cur.All(nil, &users); err != nil {
		return fmt.Errorf("loadtokens: %w", err)
	}
	if err := rdb.Del(ctx, "short-users").Err(); err != nil {
		return fmt.Errorf("loadtokens: %w", err)
	}
	for _, u := range users {
		b, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("loadtokens: %w", err)
		}
		if err := rdb.HSet(ctx, "short-users", u.Token, b).Err(); err != nil {
			return fmt.Errorf("loadtokens: %w", err)
		}
	}
	return nil
}

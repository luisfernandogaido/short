package model

import (
	"fmt"
	"time"
)

type Acesso struct {
	Ip        string
	Token     string
	Data      time.Time
	Path      string
	UserAgent string
}

func AcessoLoga(a Acesso) error {
	_, err := db.Collection("acessos").InsertOne(nil, a)
	if err != nil {
		return fmt.Errorf("acessologa: %w", err)
	}
	return nil
}

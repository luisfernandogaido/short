package model

import "log"

func setup() {
	if db != nil {
		return
	}
	if err := Ini("mongodb://localhost:27017"); err != nil {
		log.Fatal(err)
	}
}

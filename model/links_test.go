package model

import (
	"fmt"
	"testing"
)

func TestGenerateHash(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(generateHash())
	}
}

func TestLinkCreate(t *testing.T) {
	setup()
	u := User{
		Name:  "gaido",
		Token: "umtok",
	}
	link, err := LinkCreate("https://google.com", u)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(link.Destination, link.Hash)
}

func TestLinkGet(t *testing.T) {
	setup()
	link, err := LinkGet("uQYWNJG")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(link.Destination)
}

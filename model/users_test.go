package model

import (
	"fmt"
	"log"
	"testing"
)

func TestUser_Save(t *testing.T) {
	setup()
	u := User{
		Name: "Gaido",
	}
	if err := u.Save(); err != nil {
		t.Fatal(err)
	}
	fmt.Println(u.Token)
}

func TestNewUser(t *testing.T) {
	setup()
	u, err := NewUser("604ad72081afe2d7a4335e0c")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u)
}

func TestUser_Save2(t *testing.T) {
	setup()
	u, err := NewUser("604ad72081afe2d7a4335e0c")
	if err != nil {
		t.Fatal(err)
	}
	u.Name = "Maia"
	if err := u.Save(); err != nil {
		log.Fatal(err)
	}
}

func TestUser_RegenerateToken(t *testing.T) {
	setup()
	u, err := NewUser("604ad72081afe2d7a4335e0c")
	if err != nil {
		t.Fatal(err)
	}
	if err := u.RegenerateToken(); err != nil {
		t.Fatal(err)
	}
}

func TestUser_Delete(t *testing.T) {
	setup()
	u, err := NewUser("604ada317b20db652c543ad8")
	if err != nil {
		t.Fatal(err)
	}
	if err := u.Delete(); err != nil {
		t.Fatal(err)
	}
}

func TestUsers(t *testing.T) {
	setup()
	users, err := Users()
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range users {
		fmt.Println(u.Token, u.Name)
	}
}

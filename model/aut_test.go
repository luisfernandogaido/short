package model

import (
	"fmt"
	"testing"
)

func TestUserToken(t *testing.T) {
	setup()
	usr, err := UserToken("de095e41933824cfb389120dd265fc99")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(usr.Name)
}

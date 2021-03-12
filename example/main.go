package main

import (
	"log"

	"github.com/luisfernandogaido/short/api"
)

func main() {
	if err := api.Serve(":4018", "mongodb://localhost:27017", "MDKTQuSFiUleok6GA0MF"); err != nil {
		log.Fatal(err)
	}
}

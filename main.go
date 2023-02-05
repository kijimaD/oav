package main

import (
	"log"
	"os"

	"github.com/kijimaD/oav/oa"
)

func main() {
	log.SetFlags(0)

	file, err := os.Open("openapi.yml")
	if err != nil {
		panic(err)
	}
	c := oa.New(os.Stdout, file)
	err = c.Run("/pets")

	if err != nil {
		log.Fatalf("!! %+v", err)
	}
}

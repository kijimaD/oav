package main

import (
	"log"
	"os"

	"github.com/kijimaD/oav/oa"
)

func main() {
	log.SetFlags(0)

	c := oa.New(os.Stdout)
	err := c.Run("/pets")

	if err != nil {
		log.Fatalf("!! %+v", err)
	}
}

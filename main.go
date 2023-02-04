package main

import (
	"log"

	"github.com/kijimaD/oav/oa"
)

func main() {
	log.SetFlags(0)
	if err := oa.Run("/pets"); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

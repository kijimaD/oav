package main

import (
	"log"
	"net/url"
	"os"

	"github.com/kijimaD/oav/oa"
)

func main() {
	var rawURL string
	var schemaPath string
	if len(os.Args) < 2 {
		panic("not enough arguments!")
	}
	schemaPath = os.Args[1]
	rawURL = os.Args[2]

	file, err := os.Open(schemaPath)
	if err != nil {
		panic(err)
	}

	baseURL, err := url.Parse(rawURL)

	if err != nil {
		panic(err)
	}

	c := oa.New(os.Stdout, file, *baseURL)
	err = c.Run("/pets")

	if err != nil {
		log.Fatalf("!! %+v", err)
	}
}

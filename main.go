package main

import (
	"log"
	"net/url"
	"os"

	"github.com/kijimaD/oav/oa"
)

func main() {
	var rawURL string
	if len(os.Args) < 1 {
		panic("not enough arguments!")
	}
	rawURL = os.Args[1]

	file, err := os.Open("openapi.yml")
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

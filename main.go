package main

import (
	"net/url"
	"os"

	"github.com/kijimaD/oav/oa"
)

func main() {
	var schemaPath string
	if len(os.Args) < 1 {
		panic("not enough arguments!")
	}
	schemaPath = os.Args[1]

	file, err := os.Open(schemaPath)
	if err != nil {
		panic(err)
	}

	baseURL, err := url.Parse("dummy")
	if err != nil {
		panic(err)
	}

	c := oa.New(os.Stdout, file, *baseURL)
	c.DumpRoutes()
}

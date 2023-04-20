package main

import (
	"bytes"
	"io"
	"net/url"
	"os"

	"github.com/kijimaD/oav/oa"
)

// show all routes
// go run . openapi.yml
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
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		panic(err)
	}

	baseURL, err := url.Parse("dummy")
	if err != nil {
		panic(err)
	}

	c := oa.New(os.Stdout, buf, *baseURL)
	err = c.DumpRoutes()
	if err != nil {
		panic(err)
	}
}

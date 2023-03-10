package oa

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const basePath = "http://localhost:8089"

func testServer(t *testing.T) *httptest.Server {
	t.Helper()

	url, _ := url.Parse(basePath)
	l, err := net.Listen("tcp", url.Host)
	if err != nil {
		log.Fatal(err)
	}
	ts := httptest.NewUnstartedServer(routes())
	ts.Listener.Close()
	ts.Listener = l
	ts.Start()
	return ts
}

func TestValidate(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()

	// need to fix server address
	url, _ := url.Parse(basePath)

	buffer := bytes.Buffer{}
	cli := New(&buffer, strings.NewReader(schemafileA), *url)
	err := cli.Run("/pets")
	if err != nil {
		fmt.Fprint(&buffer, err)
	}

	got := buffer.String()
	assert.Contains(t, got, "GET")
	assert.Contains(t, got, "request is ok")
	assert.Contains(t, got, `"pets"`)
	assert.Contains(t, got, `"animal"`)
	assert.Contains(t, got, `"cat"`)
	assert.Contains(t, got, `"dog"`)
	assert.Contains(t, got, "response is ok")
}

func TestValidateRouteNotMatch(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()

	// need to fix server address
	url, _ := url.Parse(basePath)

	buffer := bytes.Buffer{}
	cli := New(&buffer, strings.NewReader(schemafileA), *url)
	err := cli.Run("/not_exists")
	if err != nil {
		fmt.Fprint(&buffer, err)
		got := buffer.String()
		assert.Contains(t, got, "/not_exists")
	} else {
		t.Errorf("expected: error, actual: no error")
	}

}

func TestValidateRouteInvalidResponse(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()

	// need to fix server address
	url, _ := url.Parse(basePath)

	buffer := bytes.Buffer{}
	cli := New(&buffer, strings.NewReader(schemafileB), *url)
	err := cli.Run("/pets")
	if err != nil {
		fmt.Fprint(&buffer, err)
		got := buffer.String()
		assert.Contains(t, got, "/pets")
		assert.Contains(t, got, "value must be a string")
		assert.Contains(t, got, "description\": \"pet ID")
		assert.Contains(t, got, "type\": \"string")
	} else {
		t.Errorf("expected: error, actual: no error")
	}

}

func TestDumpRoutes(t *testing.T) {
	buffer := bytes.Buffer{}
	url, _ := url.Parse(basePath)
	cli := New(&buffer, strings.NewReader(schemafileA), *url)
	err := cli.dumpRoutes()
	if err != nil {
		fmt.Fprint(&buffer, err)
	}

	got := buffer.String()
	assert.Contains(t, got, "/pets")
	assert.Contains(t, got, "Get")
	assert.Contains(t, got, "list_pets")
}

const schemafileA = `---
openapi: "3.1.0"

info:
  description: |
    ## develop
    hello world
      - list
        - A
        - B

  version: 1.0.0
  title: API Docs
  contact:
    name: kijimad
    email: norimaking777@gmail.com

servers:
  - url: http://localhost:8089
    description: go server
  - url: http://localhost:6969
    description: mock(Prism)

tags:
  - name: Pet
    description: |
      pet

paths:
  /pets:
    get:
      summary: list pets
      description: list pets
      operationId: list_pets
      tags:
        - Pet
      parameters:
        - $ref: "#/components/parameters/Limit"
      responses:
        '200':
          description: success
          content:
            application/json:
              schema:
                required:
                  - pets
                properties:
                  pets:
                    $ref: "#/components/schemas/Pets"
              examples:
                case1:
                  $ref: "#/components/examples/PetsResponse"
        default:
          description: Unexpected error
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Error"

components:
  schemas:
    Pets:
      type: array
      description: list pets
      items:
        required:
          - id
          - name
          - animal
        properties:
          id:
            type: integer
            description: pet ID
          name:
            type: string
          animal:
            type: object
            required:
              - dog
              - cat
            properties:
              dog:
                type: string
              cat:
                type: string
    Error:
      required:
        - code
        - message
      properties:
        code:
          type: integer
          format: int32
        message:
          type: string
  examples:
    PetsResponse:
      description: pets
      value:
        pets:
          - id: 1
            name: dog
            animal:
              dog: pochi
              cat: tama
          - id: 2
            name: cat
            animal:
              dog: pochi
              cat: tama
  parameters:
    Limit:
      name: limit
      in: query
      description: data count
      required: false
      schema:
        type: integer
        format: int32
`

const schemafileB = `---
openapi: "3.1.0"

info:
  description: |
    ## develop
    hello world
      - list
        - A
        - B

  version: 1.0.0
  title: API Docs
  contact:
    name: kijimad
    email: norimaking777@gmail.com

servers:
  - url: http://localhost:8089
    description: go server
  - url: http://localhost:6969
    description: mock(Prism)

tags:
  - name: Pet
    description: |
      pet

paths:
  /pets:
    get:
      summary: list pets
      description: list pets
      operationId: list_pets
      tags:
        - Pet
      parameters:
        - $ref: "#/components/parameters/Limit"
      responses:
        '200':
          description: success
          content:
            application/json:
              schema:
                required:
                  - pets
                properties:
                  pets:
                    $ref: "#/components/schemas/Pets"
              examples:
                case1:
                  $ref: "#/components/examples/PetsResponse"
        default:
          description: Unexpected error
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Error"

components:
  schemas:
    Pets:
      type: array
      description: list pets
      items:
        properties:
          id:
            type: string
            description: pet ID
    Error:
      required:
        - code
        - message
      properties:
        code:
          type: integer
          format: int32
        message:
          type: string
  examples:
    PetsResponse:
      description: pets
      value:
        pets:
          - id: "1"
            name: dog
          - id: "2"
            name: cat
  parameters:
    Limit:
      name: limit
      in: query
      description: data count
      required: false
      schema:
        type: integer
        format: int32
`

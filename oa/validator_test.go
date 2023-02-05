package oa

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

const basePath = "localhost:8080"

func TestServer(t *testing.T) {
	// need to fix server address
	l, err := net.Listen("tcp", basePath)
	if err != nil {
		log.Fatal(err)
	}
	ts := httptest.NewUnstartedServer(routes())
	ts.Listener.Close()
	ts.Listener = l
	ts.Start()
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal("status:", resp.StatusCode)
	}

	buffer := bytes.Buffer{}
	cli := New(&buffer)
	err = cli.Run("/pets")

	got := buffer.String()

	assert.Contains(t, got, "GET")
	assert.Contains(t, got, "request is ok")
	assert.Contains(t, got, `{"pets":[{"id":1},{"id":2}]}`)
	assert.Contains(t, got, "request is ok")

	if err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func TestDumpRoutes(t *testing.T) {
	buffer := bytes.Buffer{}
	cli := New(&buffer)
	doc, err := openapi3.NewLoader().LoadFromData(spec)
	if err != nil {
		t.Fatal("Failed")
	}
	cli.dumpRoutes(doc)

	got := buffer.String()
	assert.Contains(t, got, "/pets")
	assert.Contains(t, got, "Get")
	assert.Contains(t, got, "list_pets")
}

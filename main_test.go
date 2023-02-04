package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
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

	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

package main

import (
	"log"
	"testing"
)

func TestOA(t *testing.T) {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

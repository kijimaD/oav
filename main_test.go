package main

import (
	"log"
	"testing"
)

func TestOA(t *testing.T) {
	ch1 := make(chan bool)
	go exec(ch1)

	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}

	<-ch1
}

package oa

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Res struct {
	Pets []Pet `json:"pets"`
}

type Pet struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Animal Animal `json:"animal"`
}

type Animal struct {
	Dog        string     `json:"dog"`
	Cat        string     `json:"cat"`
	AnimalNest AnimalNest `json:"animalnest"`
}

type AnimalNest struct {
	ID int `json:"id"`
}

func routes() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.HandleFunc("/", root)
	mux.HandleFunc("/pets", pets)

	return
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "root")
}

func pets(w http.ResponseWriter, r *http.Request) {
	resp := Res{
		Pets: []Pet{
			{
				ID:   1,
				Name: "a",
				Animal: Animal{
					Dog: "d",
					Cat: "c",
					AnimalNest: AnimalNest{
						ID: 1,
					},
				},
			},
			{
				ID:   2,
				Name: "b",
				Animal: Animal{
					Dog: "d",
					Cat: "c",
					AnimalNest: AnimalNest{
						ID: 1,
					},
				},
			},
		},
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		panic(err)
	}
}

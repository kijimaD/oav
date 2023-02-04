package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Res struct {
	Pets []Pet `json:"pets"`
}

type Pet struct {
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
				ID: 1,
			},
			{
				ID: 2,
			},
		},
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Res struct {
	Pets []Pet `json:"pets"`
}

type Pet struct {
	ID int `json:"id"`
}

func exec(ch chan bool) {
	http.HandleFunc("/", root)
	http.HandleFunc("/pets", pets)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		panic(err)
	}

	ch <- true
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

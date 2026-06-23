package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /queue", put)
	mux.HandleFunc("GET /queue", get)

	http.ListenAndServe(":8080", mux)
}

func put(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
}

func get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
}

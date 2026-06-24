package main

import (
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)

	http.ListenAndServe(":8080", mux)
}

func handler(w http.ResponseWriter, r *http.Request) {
	qname := strings.Split(strings.Trim(r.URL.Path, "/"), "/")[0]
	// qname := r.URL.Path[1:]
	// if qname == "" {
	// http.NotFound(w, r)
	// 	return
	// }

	switch r.Method {
	case http.MethodPut:
		put(w, r, qname)
	case http.MethodGet:
		get(w, r, qname)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func put(w http.ResponseWriter, r *http.Request, qname string) {
	msg := r.URL.Query().Get("v")
	if msg == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "")
}

func get(w http.ResponseWriter, r *http.Request, qname string) {
	timeoutStr := r.URL.Query().Get("timeout")

	fmt.Fprint(w, "")
}

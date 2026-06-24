package main

import (
	"net/http"
var (
	mu     sync.Mutex
	queues = make(map[string]*queue)
)

type queue struct {
	messages []string
	waiters  []chan string
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)

	http.ListenAndServe(":8080", mux)
}

func handler(w http.ResponseWriter, r *http.Request) {
	qname := r.URL.Path[1:]

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

	mu.Lock()
	q := getQueue(qname)

	if len(q.waiters) > 0 {
		wt := q.waiters[0]
		q.waiters = q.waiters[1:]
		mu.Unlock()

		wt <- msg
	} else {
		q.messages = append(q.messages, msg)
		mu.Unlock()
	}

	w.WriteHeader(http.StatusOK)

}

func get(w http.ResponseWriter, r *http.Request, qname string) {
	timeoutStr := r.URL.Query().Get("timeout")

func getQueue(name string) *queue {
	q := queues[name]
	if q == nil {
		q = &queue{}
		queues[name] = q
	}
	return q
}

package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	mu     sync.Mutex
	queues = make(map[string]*queue)
)

type queue struct {
	messages []string
	waiters  []chan string
}

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	if *port == 0 {
		fmt.Println("usage: queue_broker -port=N")
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	addr := fmt.Sprintf(":%d", *port)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	qname := r.URL.Path[1:]
	if qname == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodPut:
		msg := r.URL.Query().Get("v")
		if msg == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		put(qname, msg)

		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		var timeout time.Duration

		timeoutStr := r.URL.Query().Get("timeout")
		if timeoutStr != "" {
			t, err := strconv.Atoi(timeoutStr)
			if err != nil || t <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			timeout = time.Duration(t) * time.Second
		}

		msg, ok := get(qname, timeout)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte(msg))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func put(qname, msg string) {
	mu.Lock()

	q := getQueue(qname)

	if len(q.waiters) > 0 {
		wt := q.waiters[0]
		q.waiters = q.waiters[1:]
		mu.Unlock()
		wt <- msg
		return
	}

	q.messages = append(q.messages, msg)
	mu.Unlock()
}

func get(qname string, timeout time.Duration) (string, bool) {
	mu.Lock()

	q := getQueue(qname)

	if len(q.messages) > 0 {
		msg := q.messages[0]
		q.messages = q.messages[1:]
		mu.Unlock()
		return msg, true
	}

	if timeout == 0 {
		mu.Unlock()
		return "", false
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	wt := make(chan string, 1)
	q.waiters = append(q.waiters, wt)
	mu.Unlock()

	select {
	case msg := <-wt:
		return msg, true

	case <-timer.C:
		mu.Lock()
		q = getQueue(qname)
		removed := false

		for i, wtr := range q.waiters {
			if wtr == wt {
				q.waiters = append(q.waiters[:i], q.waiters[i+1:]...)
				removed = true
				break
			}
		}

		mu.Unlock()

		if removed {
			return "", false
		}

		msg := <-wt
		return msg, true
	}
}

func getQueue(name string) *queue {
	q := queues[name]
	if q == nil {
		q = &queue{}
		queues[name] = q
	}
	return q
}

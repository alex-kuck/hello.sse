package main

import (
	"fmt"
	"github.com/oklog/ulid/v2"
	"net/http"
	"strconv"
	"time"
)

func main() {
	mux := http.NewServeMux()

	c := make(chan int, 1)

	go func() {
		counter := 1
		for {
			c <- counter
			counter++
			time.Sleep(1 * time.Second)
		}
	}()
	bc := NewBroadcaster(c)

	mux.HandleFunc("/events", eventHandler(bc))
	mux.Handle("/", http.FileServer(http.Dir("static")))

	http.ListenAndServe(":8080", mux)
}

func eventHandler(broadcaster *Broadcaster[int]) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		writer.Header().Set("Content-Type", "text/event-stream")
		writer.Header().Set("Cache-Control", "no-cache")
		writer.Header().Set("Connection", "keep-alive")
		c := broadcaster.Subscribe()

		go func() {
			<-request.Context().Done()
			broadcaster.Unsubscribe(c)
		}()

		for counter := range c {
			if counter%2 == 0 {
				Event{Type: "ping", Data: strconv.Itoa(counter), Id: ulid.Make().String()}.
					Publish(writer)
			}
			Event{Data: fmt.Sprintf(`{"event": %d, "key": "%s"}`, counter, ulid.Make())}.
				Publish(writer)
			writer.(http.Flusher).Flush()
		}
	}
}

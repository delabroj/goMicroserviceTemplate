package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/status", status)
	http.Handle("/time", message(time.Now().Format(time.RFC1123)))
	http.Handle("/status-log", logRequest(http.HandlerFunc(status)))
	http.Handle("/time-log", logRequest(message(time.Now().Format(time.RFC1123))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

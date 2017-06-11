package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/status", status)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

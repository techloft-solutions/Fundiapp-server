package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", hello)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

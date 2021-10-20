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
	http.HandleFunc("/", hello)

	log.Fatal(http.ListenAndServe(port, nil))
}

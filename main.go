package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func main() {

	http.HandleFunc("/", hello)

	log.Fatal(http.ListenAndServe(":80", nil))
}

package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("/home/n0ko/index.html"))
	log.Fatal(http.ListenAndServe(":9000", fs))
}

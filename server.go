package main

import (
	"fmt"
	"io"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("traffic hit /\n")
	io.WriteString(w, "Working Server")
}

func main() {
	//next iteration handle errors with the appropriate ResponseWriter
	http.HandleFunc("/", index)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("500 Not available")
	}
}

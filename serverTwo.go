package main

import (
	"fmt"
	"net/http"
)

type HttpHandler struct{}

func (h HttpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	data := []byte("n0ko's World!")
	res.Write(data)
	fmt.Printf("traffic hit\n")
}

func main() {
	handler := HttpHandler{}
	http.ListenAndServe(":8000", handler)
}

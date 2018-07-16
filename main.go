package main

import (
	"log"
	"net/http"
)

var g = NewGraph()

func main() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/badfile", badfile)
	http.HandleFunc("/path", path)
	http.HandleFunc("/", index)
	log.Println("Service starting")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
